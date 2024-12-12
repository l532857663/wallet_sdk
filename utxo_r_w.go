package wallet_sdk

import (
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/wire"
	"github.com/shopspring/decimal"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"wallet_sdk/client"
	"wallet_sdk/elastic"
	"wallet_sdk/global"
	"wallet_sdk/utils/dir"
	"wallet_sdk/utils/logutils"
)

type GetUtxoInfo struct {
	BlockHeight int64
	Client      NodeService
	Wg          sync.WaitGroup
	Address     string
}

var (
	srv *GetUtxoInfo

	SpendUTXOMap   sync.Map
	UnSpendUTXOMap sync.Map
	TxInMap        sync.Map

	AddrUTXO   map[string][]interface{}
	AddrAmount map[string]decimal.Decimal
)

func NewGetUtxoInfo(address string) {
	// 节点信息
	chainName := BTC_RegTest
	// 链接节点
	cli, err := NewNodeService(chainName)
	if err != nil {
		log.Panicln("NewNodeService", chainName)
		return
	}
	// 创建地址目录
	global.UtxoUserUnSpendPath = global.UtxoUnSpendPath + "/" + address
	global.UtxoUserSpendPath = global.UtxoSpendPath + "/" + address
	pathList := []string{
		global.UtxoUserUnSpendPath,
		global.UtxoUserSpendPath,
	}
	dir.CreateDir(pathList...)
	srv = &GetUtxoInfo{
		Client:  cli,
		Address: address,
	}
}

func GetTransferByBlockHeight(startHeight, newHigh int64) {
	logutils.LogInfof(global.LOG, "[GetTransferByBlockHeight] Start startHeight: %v, newHigh: %v", startHeight, newHigh)
	AddrUTXO = make(map[string][]interface{})
	for i := startHeight; i <= newHigh; i++ {
		srv.BlockHeight = i
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
	if txInfoLength >= 2 {
		logutils.LogInfof(global.LOG, "Get block info, block height: [%v], have tx: [%v], time: [%v]", height, txInfoLength, endTime-startTime)
	}
	for _, txInfo := range blockInfo.Transactions {
		// 添加计数器
		srv.Wg.Add(1)
		go GetUTXOInfoByTransferInfo(txInfo)
	}
	srv.Wg.Wait()
	// 处理同块前后交易
	DealTxInByBlock()
	// 处理输出
	DealTxOutByBlock()
	// 处理同块中使用的输出
	DealSpendUTXOByBlock()
}

func GetUTXOInfoByTransferInfo(txInfo *wire.MsgTx) {
	//fmt.Printf("txInfo: %+v\n", txInfo)
	defer srv.Wg.Done()
	for _, txIn := range txInfo.TxIn {
		if txIn.PreviousOutPoint.Hash.String() == global.CoinbaseHash {
			continue
		}
		key := txIn.PreviousOutPoint.String()
		TxInMap.Store(key, txIn)
	}
	txHash := txInfo.TxHash().String()
	for i, txOut := range txInfo.TxOut {
		key := fmt.Sprintf("%s:%d", txHash, i)
		UnSpendUTXOMap.Store(key, txOut)
	}
}
func DealTxInByBlock() {
	count := 0
	TxInMap.Range(func(key, value interface{}) bool {
		count++
		defer TxInMap.Delete(key)
		txIn := key.(string)
		// 判断是否有输入在当前区块生成
		if usu, ok := UnSpendUTXOMap.Load(key); ok {
			// 存在的话这个输出不需要保存
			logutils.LogInfof(global.LOG, "Spent current txIn: %v", txIn)
			SpendUTXOMap.Store(key, usu)
			UnSpendUTXOMap.Delete(key)
			return true
		}
		// 查询历史UTXO,把对应UTXO状态置为失效
		utxoPath := fmt.Sprintf("%s/%s", global.UtxoUserUnSpendPath, txIn)
		spendUtxoPath := fmt.Sprintf("%s/%s", global.UtxoUserSpendPath, txIn)
		ok, err := dir.FileExists(utxoPath)
		if err != nil {
			logutils.LogErrorf(global.LOG, "查询UTXO错误：%+v", err)
			return false
		}
		if !ok {
			return true
		}
		logutils.LogInfof(global.LOG, "Spent history txIn: %v", txIn)
		if err := os.Rename(utxoPath, spendUtxoPath); err != nil {
			logutils.LogErrorf(global.LOG, "移动UTXO错误：%+v", err)
			return false
		}
		return true
	})
}

func DealTxOutByBlock() {
	count := 0
	UnSpendUTXOMap.Range(func(key, value interface{}) bool {
		defer UnSpendUTXOMap.Delete(key)
		UTXOInfo := GetUTXOInfoForAddress(key, value)
		if UTXOInfo == nil {
			return true
		}
		// 保存UTXO到文件
		utxoName := fmt.Sprintf("%s/%s", global.UtxoUserUnSpendPath, key)
		fmt.Printf("utxoName: %+v\n", utxoName)
		SaveUTXO2File(utxoName, *UTXOInfo)
		count++
		return true
	})
	if count > 0 {
		logutils.LogInfof(global.LOG, "Save unSpend utxo len %v, From Height(%v)", count, srv.BlockHeight)
	}
}

func DealSpendUTXOByBlock() {
	count := 0
	SpendUTXOMap.Range(func(key, value interface{}) bool {
		defer SpendUTXOMap.Delete(key)
		UTXOInfo := GetUTXOInfoForAddress(key, value)
		if UTXOInfo == nil {
			return true
		}
		// 保存UTXO到文件
		utxoName := fmt.Sprintf("%s/%s", global.UtxoUserSpendPath, key)
		SaveUTXO2File(utxoName, *UTXOInfo)
		count++
		return true
	})
	if count > 0 {
		logutils.LogInfof(global.LOG, "Save spend utxo len %v, From Height(%v)", count, srv.BlockHeight)
	}
}

func GetUTXOInfoForAddress(key, value interface{}) *elastic.UnSpentsUTXO {
	addr, pkScript, txId, vout, amount := GetUTXOInfoByTxOut(key, value)
	if addr == "" {
		return nil
	}
	if addr != srv.Address {
		return nil
	}
	// 处理当前地址的UTXO
	// UTXOInfo
	UTXOInfo := &elastic.UnSpentsUTXO{
		TxId:         txId,
		Vout:         vout,
		ScriptPubKey: pkScript,
		Amount:       amount,
		Height:       srv.BlockHeight,
	}
	return UTXOInfo
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

func SaveUTXO2File(utxoName string, UTXOInfo elastic.UnSpentsUTXO) {
	// 保存到文件
	// 将结构体编码为JSON
	jsonData, err := json.Marshal(UTXOInfo)
	if err != nil {
		logutils.LogErrorf(global.LOG, "Error marshaling JSON:%v", err)
		return
	}
	// 将JSON数据写入文件
	file, err := os.Create(utxoName)
	if err != nil {
		logutils.LogErrorf(global.LOG, "Error creating file:%v", err)
		return
	}
	defer file.Close()
	_, err = file.Write(jsonData)
	if err != nil {
		logutils.LogErrorf(global.LOG, "Error writing file:%v", err)
		return
	}
}
