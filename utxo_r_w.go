package wallet_sdk

import (
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
	Stop        bool
	StopOK      bool
}

var (
	SpendUTXOMap   sync.Map
	UnSpendUTXOMap sync.Map
	TxInMap        sync.Map

	AddrAmount map[string]decimal.Decimal
)

func NewGetUtxoInfo(address string) *GetUtxoInfo {
	// 节点信息
	chainName := global.ChainName
	// 链接节点
	cli, err := NewNodeService(chainName)
	if err != nil {
		log.Panicln("NewNodeService", chainName)
		return nil
	}
	// 创建地址目录
	global.UtxoBlockHeightByUser = global.UtxoBlockHeightPath + "/" + address
	global.UtxoUserUnSpendPath = global.UtxoUnSpendPath + "/" + address
	global.UtxoUserSpendPath = global.UtxoSpendPath + "/" + address
	pathList := []string{
		global.UtxoUserUnSpendPath,
		global.UtxoUserSpendPath,
	}
	if err := dir.CreateDir(pathList...); err != nil {
		logutils.LogErrorf(global.LOG, "Error creating file:%v", err)
		return nil
	}
	return &GetUtxoInfo{
		Client:  cli,
		Address: address,
	}
}

func (srv *GetUtxoInfo) GetTransferByBlockHeight(startHeight, newHigh int64) {
	logutils.LogInfof(global.LOG, "[GetTransferByBlockHeight] Start startHeight: %v, newHigh: %v", startHeight, newHigh)
	for i := startHeight; i <= newHigh; i++ {
		if srv.Stop {
			logutils.LogInfof(global.LOG, "Stop get utxo by address: %v, block height: %+v", srv.Address, i-1)
			break
		}
		srv.BlockHeight = i
		// 检索处理UTXO
		srv.GetTransferByBlock(i)
	}
	logutils.LogInfof(global.LOG, "Stop get utxo by address: %v, block height: %+v", srv.Address, srv.BlockHeight)
	// 保存块高
	dir.SaveFile(global.UtxoBlockHeightByUser, srv.BlockHeight)
	srv.StopOK = true
}

// GetTransferByBlock 扫块获取交易数据
func (srv *GetUtxoInfo) GetTransferByBlock(height int64) {
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
		go srv.GetUTXOInfoByTransferInfo(txInfo)
	}
	srv.Wg.Wait()
	// 处理同块前后交易
	srv.DealTxInByBlock()
	// 处理输出
	srv.DealTxOutByBlock()
	// 处理同块中使用的输出
	srv.DealSpendUTXOByBlock()
}

func (srv *GetUtxoInfo) GetUTXOInfoByTransferInfo(txInfo *wire.MsgTx) {
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
func (srv *GetUtxoInfo) DealTxInByBlock() {
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

func (srv *GetUtxoInfo) DealTxOutByBlock() {
	count := 0
	UnSpendUTXOMap.Range(func(key, value interface{}) bool {
		defer UnSpendUTXOMap.Delete(key)
		UTXOInfo := srv.GetUTXOInfoForAddress(key, value)
		if UTXOInfo == nil {
			return true
		}
		// 保存UTXO到文件
		utxoName := fmt.Sprintf("%s/%s", global.UtxoUserUnSpendPath, key)
		fmt.Printf("utxoName: %+v\n", utxoName)
		dir.SaveFile(utxoName, *UTXOInfo)
		count++
		return true
	})
	if count > 0 {
		logutils.LogInfof(global.LOG, "Save unSpend utxo len %v, From Height(%v)", count, srv.BlockHeight)
	}
}

func (srv *GetUtxoInfo) DealSpendUTXOByBlock() {
	count := 0
	SpendUTXOMap.Range(func(key, value interface{}) bool {
		defer SpendUTXOMap.Delete(key)
		UTXOInfo := srv.GetUTXOInfoForAddress(key, value)
		if UTXOInfo == nil {
			return true
		}
		// 保存UTXO到文件
		utxoName := fmt.Sprintf("%s/%s", global.UtxoUserSpendPath, key)
		dir.SaveFile(utxoName, *UTXOInfo)
		count++
		return true
	})
	if count > 0 {
		logutils.LogInfof(global.LOG, "Save spend utxo len %v, From Height(%v)", count, srv.BlockHeight)
	}
}

func (srv *GetUtxoInfo) GetUTXOInfoForAddress(key, value interface{}) *elastic.UnSpentsUTXO {
	addr, pkScript, txId, vout, amount := srv.GetUTXOInfoByTxOut(key, value)
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

func (srv *GetUtxoInfo) GetUTXOInfoByTxOut(key, value interface{}) (string, string, string, int64, decimal.Decimal) {
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

func (srv *GetUtxoInfo) GetUserHeightByAddress() int64 {
	userHeight, err := dir.GetFileContent(global.UtxoBlockHeightByUser)
	if err != nil {
		logutils.LogInfof(global.LOG, "Not get user history sync height: %+v", err)
		return 0
	}
	fmt.Printf("Get userHeight: %+v\n", string(userHeight))
	start, err := strconv.ParseInt(string(userHeight), 0, 64)
	if err != nil {
		logutils.LogErrorf(global.LOG, "Get user height error: %+v", err)
		return 0
	}
	return start
}
