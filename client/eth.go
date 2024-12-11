package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
	"strconv"
	"time"
)

type EthClient struct {
	*ethclient.Client
	RpcClient *rpc.Client
	Node      *Node
}

// 以太坊系节点通用

func NewEthClient(conf *Node) (*EthClient, error) {
	var url string
	if conf.Port > 0 {
		url = fmt.Sprintf("http://%s:%d", conf.Ip, conf.Port)
	} else {
		url = conf.Ip
	}
	client, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	node := &EthClient{
		Client:    ethclient.NewClient(client),
		RpcClient: client,
		Node:      conf,
	}
	return node, nil
}

// 合约调用方法

func (c *EthClient) ContractCall(contractAddr, params, from, to string, args ...interface{}) ([]byte, error) {
	var res []byte
	// ERC20 通用ABI调用
	if len(contractAddr) != 42 {
		return res, fmt.Errorf("invalid contract address length %s", contractAddr)
	}
	data, err := Erc20Abi.Pack(params, args...)
	if err != nil {
		return res, err
	}
	result, err := c.EvmCall(from, to, data)
	if err != nil {
		return res, err
	}
	return result, nil
}
func (c *EthClient) EvmCall(fromAddr, toAddr string, data []byte) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	from := EthAddressChange(fromAddr)
	to := EthAddressChange(toAddr)
	msg := ethereum.CallMsg{
		From:     from,
		To:       &to,
		Gas:      uint64(0),
		GasPrice: big.NewInt(0),
		Value:    big.NewInt(0),
		Data:     data,
	}
	return c.Client.CallContract(ctx, msg, nil)
}

// 构建交易

func (c *EthClient) BuildTransferInfo(fromAddr, toAddr, contract, amount, gasPrice, nonce string) (interface{}, error) {
	var (
		apiTx *types.LegacyTx
		input []byte
		// 初始化 交易ETH
		gasLimit = EthGasLimit
		to       = toAddr
	)
	rawAmount, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		return nil, paramsErr("amount")
	}
	// fmt.Printf("rawAmount: %+v\n", rawAmount)
	gas, ok := new(big.Int).SetString(gasPrice, 10)
	if !ok {
		return nil, paramsErr("gasPrice")
	}
	nonceUint, err := strconv.ParseUint(nonce, 0, 64)
	if err != nil {
		return nil, paramsErr("nonce")
	}
	// 不传合约地址则为主币交易
	if contract != "" {
		// Erc20交易
		input, err = Erc20Abi.Pack("transfer", EthAddressChange(toAddr), rawAmount)
		if nil != err {
			return nil, err
		}
		rawAmount = big.NewInt(0)
		to = contract
		gasLimit = Erc20GasLimit
	}
	apiTx, err = genTransaction(to, rawAmount, gas, nonceUint, gasLimit, input)
	if nil != err {
		return nil, err
	}
	return apiTx, nil
}

func genTransaction(toAddr string, amount, gasPrice *big.Int, nonce, gasLimit uint64, input []byte) (*types.LegacyTx, error) {
	if len(toAddr) > 0 && !common.IsHexAddress(toAddr) {
		return nil, fmt.Errorf("genTransaction %s", "invalid address")
	}
	var tx *types.LegacyTx
	to := EthAddressChange(toAddr)
	tx = &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &to,
		Value:    amount,
		Data:     input,
	}
	return tx, nil
}

func (c *EthClient) SignTransferToRaw(txObj, hexPrivateKey string) (string, error) {
	// txInfo
	apiTx := &types.LegacyTx{}
	err := json.Unmarshal([]byte(txObj), apiTx)
	if err != nil {
		return "", err
	}
	prikey, err := crypto.HexToECDSA(hexPrivateKey)
	if err != nil {
		return "", err
	}

	// signer
	chainId, err := c.ChainID()
	if err != nil {
		return "", err
	}
	signer := types.NewEIP155Signer(chainId)

	// 签名数据
	signedTx, err := types.SignNewTx(prikey, signer, apiTx)
	if err != nil {
		return "", err
	}
	data, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return "", err
	}
	signTx := hexutil.Encode(data)

	return signTx, nil
}

// SendRawTransaction sends tx to node.
func (c *EthClient) SendRawTransaction(hexTx string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	var hash string
	err := c.RpcClient.CallContext(ctx, &hash, "eth_sendRawTransaction", hexTx)
	if err != nil {
		return hash, err
	}
	return hash, nil
}

// 构建合约调用

func (c *EthClient) BuildContractInfo(contract, abiContent, gasPrice, nonce, params string, args ...interface{}) (interface{}, error) {
	var (
		apiTx *types.LegacyTx
		input []byte
		// 初始化 交易ETH
		gasLimit  = Erc20GasLimit
		rawAmount = big.NewInt(0)
	)
	gas, ok := new(big.Int).SetString(gasPrice, 10)
	if !ok {
		return nil, paramsErr("gasPrice")
	}
	nonceUint, err := strconv.ParseUint(nonce, 0, 64)
	if err != nil {
		return nil, paramsErr("nonce")
	}
	// 获取合约信息，转化参数类型
	contractAbi, argsNew, err := GetAbiAndArgs(contract, abiContent, params, args)
	if err != nil {
		return nil, paramsErr("abiContent")
	}
	// Erc20交易
	input, err = contractAbi.Pack(params, argsNew...)
	if nil != err {
		return nil, err
	}
	// gasLimit * (1+len(input)/100)
	added := uint64(1 + len(input)/100)
	gasLimit = gasLimit * added
	apiTx, err = genTransaction(contract, rawAmount, gas, nonceUint, gasLimit, input)
	if nil != err {
		return nil, err
	}
	return apiTx, nil
}

// 查询使用的节点信息

func (c *EthClient) GetNodeInfo() (*Node, error) {
	return c.Node, nil
}

// 关闭链接

func (c *EthClient) Close() {
	if c.Client != nil {
		c.Client.Close()
	}
	if c.RpcClient != nil {
		c.RpcClient.Close()
	}
}

/** 链暂不支持的方法
 */

// 查询地址是否是Taproot类型

func (c *EthClient) GetAddressIsTaproot(addr string) bool {
	return false
}

// 构建部分签名操作

func (c *EthClient) BuildPSBTransfer(ins []Input, outs []Output) (interface{}, error) {
	return nil, MethodNotSupportYet
}

// 查询地址UTXO列表

func (c *EthClient) GetAddressUTXO(addr, state string) (interface{}, error) {
	return nil, MethodNotSupportYet
}

// 构建多对多交易

func (c *EthClient) BuildTransferInfoByList(unSpendUTXOList []*UnspendUTXOList, toAddrList []*ToAddrDetail, gasPrice, changeAddr string) (interface{}, error) {
	return nil, MethodNotSupportYet
}

// 多个地址的签名出账

func (c *EthClient) SignListAndSendTransfer(txObj string, hexPrivateKeys []string) (string, error) {
	return "", MethodNotSupportYet
}

// 构建部分签名操作

func (c *EthClient) GenerateSignedListingPSBTBase64(ins *Input, outs *Output) (interface{}, error) {
	return nil, MethodNotSupportYet
}

func (c *EthClient) GetBlockHashByHeight(height int64) (string, error) {
	return "", MethodNotSupportYet
}

func (c *EthClient) GetBlockInfoByHeight(height int64) (interface{}, error) {
	return nil, MethodNotSupportYet
}

func (c *EthClient) GetBlockInfoByHash(hash string) (interface{}, error) {
	return nil, MethodNotSupportYet
}

func (c *EthClient) GetParams() interface{} {
	return nil
}
