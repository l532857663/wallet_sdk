package client

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"google.golang.org/grpc"
)

type TronClient struct {
	client *client.GrpcClient
	Node   *Node
}

func NewTronClient(conf *Node) (*TronClient, error) {
	address := conf.Ip + ":" + strconv.FormatUint(conf.Port, 10)
	client := client.NewGrpcClientWithTimeout(address, 10*time.Second)

	if err := client.Start(grpc.WithInsecure()); nil != err {
		return nil, err
	}

	ret := TronClient{
		client: client,
	}

	return &ret, nil
}

// 查询地址余额
func (t *TronClient) GetBalance(addr, state string) (*big.Int, error) {
	account, err := t.client.GetAccount(addr)
	if err != nil {
		if strings.Contains(err.Error(), "account not found") {
			return big.NewInt(0), nil
		}
		return nil, err
	}
	return big.NewInt(account.Balance), err
}

// 查询地址代币余额
func (t *TronClient) GetBalanceByContract(addr, contractAddr string) (*big.Int, error) {
	return t.client.TRC20ContractBalance(addr, contractAddr)
}

// 查询交易信息
func (t *TronClient) GetTransactionByHash(txHash string) (interface{}, bool, error) {
	txInfo, err := t.client.GetTransactionByID(txHash)
	return txInfo, false, err
}

// 查询合约精度
func (t *TronClient) GetDecimals(contractAddr string) (*big.Int, error) {
	return t.client.TRC20GetDecimals(contractAddr)
}

// 查询合约symbol
func (t *TronClient) GetSymbol(contractAddr string) (string, error) {
	return t.client.TRC20GetName(contractAddr)
}

// 查询建议手续
func (t *TronClient) SuggestGasPrice() *big.Int {
	return big.NewInt(0)
}

// 查询地址的nonce
func (t *TronClient) GetNonce(addr, param string) (uint64, error) {
	return 0, nil
}

// 查询chain_id
func (t *TronClient) ChainID() (*big.Int, error) {
	return big.NewInt(0), nil
}

// 查询最新区块高度
func (t *TronClient) GetBlockHeight() (int64, error) {
	block, err := t.client.GetNowBlock()
	if nil != err {
		return -1, err
	}

	return block.BlockHeader.RawData.Number, nil
}

// 构建交易
func (t *TronClient) BuildTransferInfo(fromAddr, toAddr, contract, amount, gasPrice, nonce string) (interface{}, error) {
	rawAmount, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		return nil, fmt.Errorf("The params error!")
	}
	// 不传合约地址则为主币交易
	if contract != "" {
		// TRC20交易
		return t.client.TRC20Send(fromAddr, toAddr, contract, rawAmount, Trc20FeeLimit)
	}
	// TRX交易
	return t.client.Transfer(fromAddr, toAddr, rawAmount.Int64())
}

func (t *TronClient) SignTransferToRaw(txObj, hexPrivateKey string) (string, error) {
	// txInfo
	apiTx := &api.TransactionExtention{}
	err := json.Unmarshal([]byte(txObj), apiTx)
	if err != nil {
		return "", err
	}
	tx := apiTx.Transaction
	if tx == nil {
		return "", fmt.Errorf("Transcation is nil")
	}
	// Prikey
	prikey, err := crypto.HexToECDSA(hexPrivateKey)
	if err != nil {
		return "", err
	}

	// 签名数据
	signature, err := crypto.Sign(apiTx.Txid, prikey)
	if err != nil {
		return "", err
	}
	apiTx.Transaction.Signature = append(apiTx.Transaction.Signature, signature)
	signTx, err := json.Marshal(apiTx)
	if err != nil {
		return "", fmt.Errorf("Marshal trx rawTx error")
	}
	return string(signTx), nil
}

// SendRawTransaction sends tx to node.
func (t *TronClient) SendRawTransaction(hexTx string) (string, error) {
	apiTx := &api.TransactionExtention{}
	err := json.Unmarshal([]byte(hexTx), apiTx)
	if err != nil {
		return "", fmt.Errorf("The tx error: %+v", err)
	}
	tx := apiTx.Transaction
	if tx == nil {
		return "", fmt.Errorf("Transcation is nil")
	}
	ret, err := t.client.Broadcast(apiTx.Transaction)
	if err != nil {
		return "", err
	}

	if !ret.Result {
		return "", fmt.Errorf("TronClient BroadTransaction fatal")
	}

	txHash := hex.EncodeToString(apiTx.GetTxid())
	return txHash, nil
}

// 查询使用的节点信息
func (t *TronClient) GetNodeInfo() (*Node, error) {
	return t.Node, nil
}

// 查询地址是否是Taproot类型
func (t *TronClient) GetAddressIsTaproot(addr string) bool {
	return false
}

// 关闭链接
func (t *TronClient) Close() {
	if t.client != nil {
		t.client.Stop()
	}
}

/** 链暂不支持的方法
 */

// 构建合约调用
func (t *TronClient) BuildContractInfo(contract, abiContent, gasPrice, nonce, params string, args ...interface{}) (interface{}, error) {
	return nil, MethodNotSupportYet
}

// 构建部分签名操作
func (t *TronClient) BuildPSBTransfer(ins []Input, outs []Output) (interface{}, error) {
	return nil, MethodNotSupportYet
}

// 查询地址UTXO列表
func (t *TronClient) GetAddressUTXO(addr, state string) (interface{}, error) {
	return nil, MethodNotSupportYet
}

// 构建多对多交易
func (t *TronClient) BuildTransferInfoByList(unSpendUTXOList []*UnspendUTXOList, toAddrList []*ToAddrDetail, gasPrice, changeAddr string) (interface{}, error) {
	return nil, MethodNotSupportYet
}

// 多个地址的签名出账
func (t *TronClient) SignListAndSendTransfer(txObj string, hexPrivateKeys []string) (string, error) {
	return "", MethodNotSupportYet
}

// 构建部分签名操作
func (t *TronClient) GenerateSignedListingPSBTBase64(ins *Input, outs *Output) (interface{}, error) {
	return nil, MethodNotSupportYet
}

// 构建合约调用
func (t *TronClient) GetContractInfoByFunc(contractAddr, funcName string, args ...interface{}) (interface{}, error) {
	return nil, MethodNotSupportYet
}

func (t *TronClient) GetBlockHashByHeight(height int64) (string, error) {
	return "", MethodNotSupportYet
}

func (t *TronClient) GetBlockInfoByHeight(height int64) (interface{}, error) {
	return nil, MethodNotSupportYet
}

func (t *TronClient) GetBlockInfoByHash(hash string) (interface{}, error) {
	return nil, MethodNotSupportYet
}

func (t *TronClient) GetParams() interface{} {
	return nil
}
