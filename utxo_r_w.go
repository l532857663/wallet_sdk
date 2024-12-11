package wallet_sdk

import (
	"fmt"
	"github.com/btcsuite/btcd/wire"
	"github.com/shopspring/decimal"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
	"wallet_sdk/client"
	"wallet_sdk/elastic"
	"wallet_sdk/global"
	"wallet_sdk/utils/logutils"
)

type GetUtxoInfo struct {
	BlockHeight int
	Client      NodeService
	Wg          sync.WaitGroup
}

var (
	srv *GetUtxoInfo

	SpendUTXOMap   sync.Map
	UnSpendUTXOMap sync.Map
	AddrUTXO       map[string][]interface{}
	AddrAmount     map[string]decimal.Decimal
)

func InitNode() {
	// 节点信息
	chainName := BTC_RegTest
	// 链接节点
	cli, err := NewNodeService(chainName)
	if err != nil {
		log.Panicln("NewNodeService", chainName)
		return
	}
	srv = &GetUtxoInfo{
		Client: cli,
	}
}

func GetTransferByBlockHeight(startHeight, newHigh int64) {
	logutils.LogInfof(global.LOG, "[GetTransferByBlockHeight] Start startHeight: %v, newHigh: %v", startHeight, newHigh)
	AddrUTXO = make(map[string][]interface{})
	for i := startHeight; i <= newHigh; i++ {
		GetTransferByBlock(i)
	}
	for key, val := range AddrUTXO {
		fmt.Printf("wch---- AddrUTXO: %v, %+v\n", key, len(val))
	}
}

func GetUTXOBacktrackFromHeight(stopHeight int64) {
	logutils.LogInfof(global.LOG, "[GetUTXOBacktrackFromHeight] Start newHigh: %v", stopHeight)
	AddrUTXO = make(map[string][]interface{})
	for i := stopHeight; i >= 0; i-- {
		GetTransferByBlock(i)
	}
	for key, val := range AddrUTXO {
		fmt.Printf("wch---- AddrUTXO: %v, %+v\n", key, len(val))
	}
}

// GetTransferByBlock 扫块获取交易数据
func GetTransferByBlock(height int64) {
	funcName := "GetTransferByBlock"
	startTime := time.Now().Unix()
	blockInfoInc, err := srv.Client.GetBlockInfoByHeight(height)
	if err != nil {
		logutils.LogErrorf(global.LOG, "[%s]GetBlockInfoByHash error: %+v", funcName, err)
		return
	}
	endTime := time.Now().Unix()
	// fmt.Printf("wch---- blockInfo: %+v\n", blockInfo)
	blockInfo := blockInfoInc.(*wire.MsgBlock)
	txInfoLength := len(blockInfo.Transactions)
	logutils.LogInfof(global.LOG, "Get block info, block height: [%v], have tx: [%v], time: [%v]", height, txInfoLength, endTime-startTime)
	for _, txInfo := range blockInfo.Transactions {
		// 添加计数器
		srv.Wg.Add(1)
		go GetUTXOInfoByTransferInfo(txInfo)
	}
	srv.Wg.Wait()
	// 处理同块前后交易
	DealTxInByBlock()
	// 处理输出
	DealTxOutByBlock(height)
}

func GetUTXOInfoByTransferInfo(txInfo *wire.MsgTx) {
	//fmt.Printf("txInfo: %+v\n", txInfo)
	defer srv.Wg.Done()
	for _, txIn := range txInfo.TxIn {
		if txIn.PreviousOutPoint.Hash.String() == global.CoinbaseHash {
			continue
		}
		key := txIn.PreviousOutPoint.String()
		SpendUTXOMap.Store(key, txIn)
	}
	txHash := txInfo.TxHash().String()
	for i, txOut := range txInfo.TxOut {
		key := fmt.Sprintf("%s:%d", txHash, i)
		UnSpendUTXOMap.Store(key, txOut)
	}
}
func DealTxInByBlock() {
	count := 0
	SpendUTXOMap.Range(func(key, value interface{}) bool {
		count++
		defer SpendUTXOMap.Delete(key)
		txIn := key.(string)
		// 判断是否有输入在当前区块生成
		if _, ok := UnSpendUTXOMap.Load(key); ok {
			// 存在的话这个输出不需要保存
			fmt.Println("Spent current txIn:", txIn)
			UnSpendUTXOMap.Delete(key)
			return true
		}
		// 查询历史UTXO,把对应UTXO状态置为失效
		fmt.Println("Spent history txIn:", txIn)
		return true
	})
}

func DealTxOutByBlock(height int64) {
	count := 0
	UnSpendUTXOMap.Range(func(key, value interface{}) bool {
		defer UnSpendUTXOMap.Delete(key)
		addr, pkScript, txId, vout, amount := GetUTXOInfoByTxOut(key, value)
		if addr == "" {
			return true
		}
		// UTXOInfo
		UTXOInfo := elastic.UnSpentsUTXO{
			TxId:         txId,
			Vout:         vout,
			ScriptPubKey: pkScript,
			Amount:       amount,
			Height:       height,
		}
		AddrUTXO[addr] = append(AddrUTXO[addr], UTXOInfo)
		count++
		return true
	})
}

func GetUTXOInfoByTxOut(key, value interface{}) (string, string, string, int64, decimal.Decimal) {
	zero := decimal.Zero
	txOut := value.(*wire.TxOut)
	// PKScript -> address, addrType
	addr, err := client.GetAddressByPKScript(txOut.PkScript, srv.Client.GetParams())
	if err != nil {
		return "", "", "", 0, zero
	}
	pkScript := fmt.Sprintf("%x", txOut.PkScript)
	// txInfo
	txInfo := strings.Split(key.(string), ":")
	if len(txInfo) != 2 {
		return "", "", "", 0, zero
	}
	// vout
	vout, err := strconv.ParseInt(txInfo[1], 0, 64)
	if err != nil {
		return "", "", "", 0, zero
	}
	// value sats
	amount := decimal.NewFromInt(txOut.Value)
	return addr, pkScript, txInfo[0], vout, amount
}
