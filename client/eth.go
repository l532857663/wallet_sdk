package client

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"
	"wallet_sdk/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
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

// 查询地址余额
func (c *EthClient) GetBalance(addr, state string) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	account := EthAddressChange(addr)
	defer cancel()
	var result hexutil.Big
	err := c.RpcClient.CallContext(ctx, &result, "eth_getBalance", account, state)
	return (*big.Int)(&result), err
}

// 查询地址代币余额
func (c *EthClient) GetBalanceByContract(addr, contractAddr string) (*big.Int, error) {
	result, err := c.ContractCall(contractAddr, "balanceOf", addr, contractAddr, EthAddressChange(addr))
	if err != nil {
		return nil, err
	}
	return utils.ByteTobigInt(result), nil
}

// 查询交易信息
func (c *EthClient) GetTransactionByHash(txHash string) (interface{}, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	return c.TransactionByHash(ctx, common.HexToHash(txHash))
}

// 查询合约精度
func (c *EthClient) GetDecimals(contractAddr string) (*big.Int, error) {
	result, err := c.ContractCall(contractAddr, "decimals", contractAddr, contractAddr)
	if err != nil {
		return nil, err
	}
	return utils.ByteTobigInt(result), nil
}

// 查询合约symbol
func (c *EthClient) GetSymbol(contractAddr string) (string, error) {
	result, err := c.ContractCall(contractAddr, "symbol", contractAddr, contractAddr)
	if err != nil {
		return "", err
	}
	// NOTE：该方法返回数据总长度为96
	if len(result) != 96 || result[31] != 0x20 {
		return "", err
	}
	res := result[64 : 64+int(result[63])]
	// fmt.Printf("byte: %+v\n", result)
	return string(res), nil
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

// SuggestGasPrice retrieves the currently suggested gas price to allow a timely
// execution of a transaction.
func (c *EthClient) SuggestGasPrice() *big.Int {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	gasPrice, err := c.Client.SuggestGasPrice(ctx)
	if err != nil {
		return big.NewInt(0)
	}

	return gasPrice
}

// 查询地址的nonce
func (c *EthClient) GetNonce(addr, param string) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	method := "eth_getTransactionCount"
	var strNonce string
	err := c.RpcClient.CallContext(ctx, &strNonce, method, addr, param)
	if err != nil {
		return 0, err
	}
	nonce := utils.HexTobigInt(strNonce)
	return nonce.Uint64(), nil
}

// 查询chain_id
func (c *EthClient) ChainID() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	id, err := c.Client.ChainID(ctx)
	if err != nil {
		return nil, err
	}
	return id, nil
}

// 查询最新区块高度
func (c *EthClient) GetBlockHeight() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	method := "eth_blockNumber"
	var strHeight string
	err := c.RpcClient.CallContext(ctx, &strHeight, method)
	if err != nil {
		return 0, err
	}
	height := utils.HexTobigInt(strHeight)
	return height.Int64(), nil
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

func (c *EthClient) SignAndSendTransfer(txObj, hexPrivateKey string, chainId *big.Int, idx int) (string, error) {
	// txInfo
	apiTx := &types.LegacyTx{}
	err := json.Unmarshal([]byte(txObj), apiTx)
	if err != nil {
		return "", err
	}
	// fmt.Printf("apiTx: %+v\n", apiTx)
	// Prikey
	fmt.Printf("prikey: %+v\n", hexPrivateKey)
	prikey, err := crypto.HexToECDSA(hexPrivateKey)
	if err != nil {
		return "", err
	}
	//--------------------------------------------------------
	// 数据hash
	tx := types.NewTx(apiTx)
	chainId = big.NewInt(1)
	//--------------------------------------------------------

	// signer
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

	//--------------------------------------------------------
	h := signer.Hash(tx).String()
	fmt.Printf("wch----- signTx: %+v\n hash: %+v\n", signTx, h)
	return "", nil
	//--------------------------------------------------------

	// 发送交易
	txHash, err := c.SendRawTransaction(signTx)
	if err != nil {
		return "", err
	}
	return txHash, nil
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
	contractAbi, argsNew, err := GetAbiAndArgs(abiContent, params, args)
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

// 查询地址是否是Taproot类型
func (c *EthClient) GetAddressIsTaproot(addr string) bool {
	return false
}

// 关闭链接
func (node *EthClient) Close() {
	if node.Client != nil {
		node.Client.Close()
	}
	if node.RpcClient != nil {
		node.RpcClient.Close()
	}
}

/** 链暂不支持的方法
 */

// 构建部分签名操作
func (c *EthClient) BuildPSBTransfer(ins []Input, outs []Output) (interface{}, error) {
	return nil, fmt.Errorf("This method is not supported yet!")
}

// 查询地址UTXO列表
func (c *EthClient) GetAddressUTXO(addr, state string) (interface{}, error) {
	return nil, fmt.Errorf("This method is not supported yet!")
}

// 构建多对多交易
func (c *EthClient) BuildTransferInfoByList(unSpendUTXOList []*UnspendUTXOList, toAddrList []*ToAddrDetail, gasPrice, changeAddr string) (interface{}, error) {
	return nil, fmt.Errorf("This method is not supported yet!")
}

// 多个地址的签名出账
func (c *EthClient) SignListAndSendTransfer(txObj string, hexPrivateKeys []string) (string, error) {
	return "", fmt.Errorf("This method is not supported yet!")
}
