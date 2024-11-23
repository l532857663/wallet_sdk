package wallet_sdk

import (
	"fmt"
	"strconv"
	"wallet_sdk/client"

	"github.com/shopspring/decimal"
)

/**
 * 查询最新块高
 *
 * Params (chainName)
 * chainName:
 *   链名称
 */
func GetBlockHeight(chainName string) *CommonResp {
	res := &CommonResp{}
	funcName := "GetBlockHeight"

	// 链接节点
	cli, err := NewNodeService(chainName)
	defer cli.Close()
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] new node client error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	height, err := cli.GetBlockHeight()
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] get block height error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// 返回结果
	res.Status = ResSuccess
	res.Data = strconv.FormatInt(height, 10)
	return res
}

/**
 * 请求nonce值
 *
 * Params (chainName, address, params string)
 * chainName:
 *   链名称
 * address:
 *   查询地址
 * params:
 *   参数："latest"、"earliest"或"pending"
 */
func GetNonce(chainName, address, params string) *CommonResp {
	res := &CommonResp{}
	funcName := "GetNonce"

	// 链接节点
	cli, err := NewNodeService(chainName)
	defer cli.Close()
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] new node client error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	nonce, err := cli.GetNonce(address, params)
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] get nonce error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// 返回结果
	res.Status = ResSuccess
	res.Data = strconv.FormatUint(nonce, 10)
	return res
}

/**
 * 请求建议GasPrice
 *
 * Params (chainName string)
 * chainName:
 *   链名称
 */
func GetGasPrice(chainName string) *TransactionGasPriceResp {
	res := &TransactionGasPriceResp{}
	funcName := "GetGetGasPrice"

	// 链接节点
	cli, err := NewNodeService(chainName)
	defer cli.Close()
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] new node client error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// BTC和ETH节点手续费保留位数不同
	nodeInfo, err := cli.GetNodeInfo()
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] get node info error: %+v", funcName, err)
		res.Status = resp
		return res
	}
	var length int32
	var baseDecimal decimal.Decimal
	zero := decimal.NewFromInt(0)
	switch ChainRelationMap[nodeInfo.ChainType] {
	case MainCoinBTC:
		length = client.SatoshiLength
		baseDecimal = client.BtcBaseDecimal
	case MainCoinETH:
		length = client.GweiLength
		baseDecimal = client.EthBaseDecimal
	default:
		length = 0
		baseDecimal = zero
	}

	gasPrice := cli.SuggestGasPrice()
	gas := decimal.NewFromBigInt(gasPrice, 0).Div(baseDecimal)
	if gas.Cmp(zero) <= 0 {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] No gas price received", funcName)
		res.Status = resp
		return res
	}
	fast := gas.Mul(GasFast).StringFixed(length)
	high := gas.Mul(GasHigh).StringFixed(length)
	average := gas.Mul(GasAverage).StringFixed(length)

	gasInfo := &GasPriceInfo{
		Fast:    fast,
		High:    high,
		Average: average,
		Low:     gas.String(),
	}

	// 返回结果
	res.Status = ResSuccess
	res.Data = gasInfo
	return res
}

/**
 * 查询地址主币余额
 *
 * Params (chainName, address string)
 * chainName:
 *   链名称
 * address:
 *   查询地址
 */
func GetBalanceByAddress(chainName, address string) *CommonResp {
	res := &CommonResp{}
	funcName := "GetBalanceByAddress"

	fmt.Printf("wch----- test1\n")
	// 链接节点
	cli, err := NewNodeService(chainName)
	defer cli.Close()
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] new node client error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	fmt.Printf("wch----- test1-1\n")
	// 请求地址余额
	balance, err := cli.GetBalance(address, StateLatest)
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] get balance error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// 返回结果
	res.Status = ResSuccess
	res.Data = balance.String()
	return res
}

/**
 * 查询地址代币合约余额
 *
 * Params (chainName, address, contract string)
 * chainName:
 *   链名称
 * address:
 *   查询地址
 * contract:
 *   合约地址
 */
func GetBalanceByAddressAndContract(chainName, address, contract string) *CommonResp {
	res := &CommonResp{}
	funcName := "GetBalanceByAddressAndContract"

	// 链接节点
	cli, err := NewNodeService(chainName)
	defer cli.Close()
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] new node client error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// 请求合约地址余额
	balance, err := cli.GetBalanceByContract(address, contract)
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] get balance error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// 返回结果
	res.Status = ResSuccess
	res.Data = balance.String()
	return res
}

/**
 * 查询交易信息
 *
 * Params (chainName, txHash string)
 * chainName:
 *   链名称
 * txHash:
 *   交易hash
 */
func GetTransaction(chainName, txHash string) *TransactionInfoResp {
	res := &TransactionInfoResp{}
	funcName := "GetTransaction"

	// 链接节点
	cli, err := NewNodeService(chainName)
	defer cli.Close()
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] new node client error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	txInfo, isPending, err := cli.GetTransactionByHash(txHash)
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] get transaction by hash error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	resp := &TransactionInfo{
		TxInfo:    txInfo,
		IsPending: isPending,
	}

	// 返回结果
	res.Status = ResSuccess
	res.Data = resp
	return res
}

/**
 * 查询合约信息
 *
 * Params (chainName, contract string)
 * chainName:
 *   链名称
 * contract:
 *   合约地址
 */
func GetContractInfo(chainName, contract string) *ContractInfoResp {
	res := &ContractInfoResp{}
	funcName := "GetContractInfo"

	// 链接节点
	cli, err := NewNodeService(chainName)
	defer cli.Close()
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] new node client error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// 请求合约地址精度
	decimals, err := cli.GetDecimals(contract)
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] get contract decimals error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// 请求合约symbol
	symbol, err := cli.GetSymbol(contract)
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] get contract symbol error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	contractInfo := &ContractInfo{
		Decimals: decimals.String(),
		Symbol:   symbol,
	}

	// 返回结果
	res.Status = ResSuccess
	res.Data = contractInfo
	return res
}

/**
 * 查询地址UTXO
 *
 * Params (chainName, address string)
 * chainName:
 *   链名称
 * address:
 *   查询地址
 */
func GetUTXOListByAddress(chainName, address string) *AddressUTXOListResp {
	res := &AddressUTXOListResp{}
	funcName := "GetUTXOListByAddress"

	// 链接节点
	cli, err := NewNodeService(chainName)
	defer cli.Close()
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] new node client error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// 请求地址余额
	utxoList, err := cli.GetAddressUTXO(address, StateLatest)
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] get balance error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// 返回结果
	res.Status = ResSuccess
	res.Data = utxoList
	return res
}

/**
 * 查询合约信息
 *
 * Params (chainName, contract string)
 * chainName:
 *   链名称
 * contract:
 *   合约地址
 */
func GetContractInfoByFunc(chainName, contract, contractFunc string, args ...interface{}) *ContractInfoResp {
	res := &ContractInfoResp{}
	funcName := "GetContractInfoByFunc"

	// 链接节点
	cli, err := NewNodeService(chainName)
	defer cli.Close()
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] new node client error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// 请求合约地址精度
	info, err := cli.GetContractInfoByFunc(contract, contractFunc, args...)
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] get contract [%s] info by funcName error: %+v", funcName, contractFunc, err)
		res.Status = resp
		return res
	}
	fmt.Printf("wch---- info: %+v\n", info)

	contractInfo := &ContractInfo{
		Decimals: "aaa",
		Symbol:   "test",
	}

	// 返回结果
	res.Status = ResSuccess
	res.Data = contractInfo
	return res
}
