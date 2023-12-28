package wallet_sdk

import (
	"encoding/json"
	"fmt"
	"wallet_sdk/client"
)

/**
 * 构建交易
 *
 * Params (chainName, fromAddr, toAddr, amount, contract, gasPrice, nonce string)
 * chainName:
 *   链名称
 * fromAddr:
 *   出账地址
 * toAddr:
 *   目标地址
 * amount:
 *   发送金额
 * contract:
 *   合约地址 发送代币的合约地址
 * gasPrice:
 *   gas price 单位：Gwei
 * nonce:
 *   发送账户的nonce
 */
func BuildTransferInfo(chainName, fromAddr, toAddr, amount, contract, gasPrice, nonce string) *CommonResp {
	res := &CommonResp{}
	funcName := "BuildTransferInfo"

	// 链接节点
	cli, err := NewNodeService(chainName)
	defer cli.Close()
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] new node client error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// Gas费单位转化
	gasPrice = client.EthToGwei(gasPrice)
	fmt.Printf("test gasPrice: %+v\n", gasPrice)

	// 创建交易结构
	TxInfo, err := cli.BuildTransferInfo(fromAddr, toAddr, contract, amount, gasPrice, nonce)
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] build transfer info error: %+v", funcName, err)
		res.Status = resp
		return res
	}
	// fmt.Printf("test Tx info: %+v\n", TxInfo)

	// 字符串
	signData, err := json.Marshal(TxInfo)
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] json.Marshal transfer info error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// 返回结果
	res.Status = ResSuccess
	res.Data = string(signData)
	return res
}

/**
 * 合约调用
 *
 * Params (chainName, contract, abiContent, gasPrice, nonce, params string, args ...interface{})
 * chainName:
 *   链名称
 * contract:
 *   合约地址
 * abiContent:
 *   合约ABI信息
 * gasPrice:
 *   gas price 单位：Gwei
 * nonce:
 *   发送账户的nonce
 * params:
 *   调用的合约方法
 * args:
 *   该合约方法的参数
 */
func BuildContractInfo(chainName, contract, abiContent, gasPrice, nonce, params string, args ...interface{}) *CommonResp {
	res := &CommonResp{}
	funcName := "BuildContractInfo"

	// 链接节点
	cli, err := NewNodeService(chainName)
	defer cli.Close()
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] new node client error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// Gas费单位转化
	gasPrice = client.EthToGwei(gasPrice)

	contractInfo, err := cli.BuildContractInfo(contract, abiContent, gasPrice, nonce, params, args...)
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] build contract info error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// 字符串
	signData, err := json.Marshal(contractInfo)
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] json.Marshal transfer info error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// 返回结果
	res.Status = ResSuccess
	res.Data = string(signData)
	return res
}

/**
 * 签名并广播交易
 *
 * Params (chainName, priKey string, apiTx interface{})
 * chainName:
 *   链名称
 * priKey:
 *   私钥
 * apiTx:
 *   交易信息
 */
func SignAndSendTransferInfo(chainName, priKey, apiTx string) *CommonResp {
	res := &CommonResp{}
	funcName := "SignTransferInfo"

	// 链接节点
	cli, err := NewNodeService(chainName)
	defer cli.Close()
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] new node client error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	chainId, _ := cli.ChainID()

	var txHash string
	if false {
		// 签名并广播交易
		txHash, err = cli.SignAndSendTransfer(apiTx, priKey, chainId, 0)
	} else {
		// 签名并广播TaprootUTXO
		txHash, err = cli.SignAndSendTaprootTransfer(apiTx, priKey, chainId, 0)
	}
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] sign transfer error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// 返回结果
	res.Status = ResSuccess
	res.Data = txHash
	return res
}

/**
 * 构建交易
 *
 * Params (chainName, fromAddr, toAddr, amount, gasPrice string)
 * chainName:
 *   链名称
 * fromAddr:
 *   出账地址
 * toAddr:
 *   目标地址
 * amount:
 *   发送金额
 * gasPrice:
 *   gas price 单位：Gwei
 */
func BuildTransferInfoByBTC(chainName, fromAddr, toAddr, amount, gasPrice string) *CommonResp {
	res := &CommonResp{}
	funcName := "BuildTransferInfoByBTC"

	// 链接节点
	cli, err := NewNodeService(chainName)
	defer cli.Close()
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] new node client error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// Gas费单位转化
	fmt.Printf("test amount: %+v\n", amount)
	fmt.Printf("test gasPrice: %+v\n", gasPrice)

	// 创建交易结构
	TxInfo, err := cli.BuildTransferInfo(fromAddr, toAddr, "", amount, gasPrice, "")
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] build transfer info error: %+v", funcName, err)
		res.Status = resp
		return res
	}
	fmt.Printf("test Tx info: %+v\n", TxInfo)

	// 字符串
	signData, err := json.Marshal(TxInfo)
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] json.Marshal transfer info error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// 返回结果
	res.Status = ResSuccess
	res.Data = string(signData)
	return res
}

func BuildPSBTransferInfo(chainName, gasPrice string) *CommonResp {
	res := &CommonResp{}
	funcName := "BuildPSBTransferInfo"

	// 链接节点
	cli, err := NewNodeService(chainName)
	defer cli.Close()
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] new node client error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// Gas费单位转化
	fmt.Printf("test gasPrice: %+v\n", gasPrice)

	// 创建交易结构
	TxInfo, err := cli.BuildPSBTransfer(Ins, Outs)
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] build transfer info error: %+v", funcName, err)
		res.Status = resp
		return res
	}
	fmt.Printf("test Tx info: %+v\n", TxInfo)

	// 字符串
	signData, err := json.Marshal(TxInfo)
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] json.Marshal transfer info error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// 返回结果
	res.Status = ResSuccess
	res.Data = string(signData)
	return res
}
