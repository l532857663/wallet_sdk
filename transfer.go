package wallet_sdk

import (
	"encoding/json"
	"fmt"
	"wallet_sdk/client"

	"github.com/shopspring/decimal"
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
	funcName := "SignAndSendTransferInfo"

	// 链接节点
	cli, err := NewNodeService(chainName)
	defer cli.Close()
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] new node client error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// 签名交易
	signTx, err := cli.SignTransferToRaw(apiTx, priKey)
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] sign transfer error: %+v", funcName, err)
		res.Status = resp
		return res
	}
	// 广播交易
	txHash, err := cli.SendRawTransaction(signTx)
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] Broadcast SendRawTransaction fatal: %+v", funcName, err)
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

func BuildPSBTransferInfo(chainName, priKey, gasPrice string) *CommonResp {
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

	// 输入输出
	In := &client.Input{
		TxId:       "f54e6280170f929fcf10630d2766125b9422e69855f2b17e6edae218e874413f",
		VOut:       0,
		Address:    "tb1pfzl0rw44mkgevdauhrtzy5kdztjezyq0rnfqfppzxtnrwzdj553qvz6lux",
		PrivateKey: priKey,
	}
	Out := &client.Output{
		Address: "tb1pfzl0rw44mkgevdauhrtzy5kdztjezyq0rnfqfppzxtnrwzdj553qvz6lux",
		Amount:  1000,
	}

	// 创建交易结构
	signData, err := cli.GenerateSignedListingPSBTBase64(In, Out)
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] build transfer info error: %+v", funcName, err)
		res.Status = resp
		return res
	}
	fmt.Printf("test Tx info: %+v\n", signData)

	// // 字符串
	// signData, err := json.Marshal(TxInfo)
	// if err != nil {
	// 	resp := ResFailed
	// 	resp.Message = fmt.Sprintf("[%s] json.Marshal transfer info error: %+v", funcName, err)
	// 	res.Status = resp
	// 	return res
	// }

	// 返回结果
	res.Status = ResSuccess
	res.Data = signData.(string)
	return res
}

/**
 * 构建多对多交易
 *
 * Params (chainName string, fromAddrs []string, vins []ChooseUTXO, toAddrs, amounts []string, gasPrice, changeAddr string)
 * chainName:
 *   链名称
 * fromAddrs:
 *   出账地址列表
 * txHashs:
 *   出账地址UTXO哈希列表
 * toAddrs:
 *   目标地址列表
 * amounts:
 *   目标金额
 * gasPrice:
 *   gas price 单位：Gwei
 */
func BuildTransferInfoByBTCList(chainName string, fromAddrs []string, vins []ChooseUTXO, toAddrs, amounts []string, gasPrice, changeAddr string) *CommonResp {
	res := &CommonResp{}
	funcName := "BuildTransferInfoByBTCList"

	// 链接节点
	cli, err := NewNodeService(chainName)
	defer cli.Close()
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] new node client error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// 最终使用UTXO
	var useUTXOList []*client.UnspendUTXOList
	// 查询FROM地址的UTXO
	for i, fromAddr := range fromAddrs {
		utxoInterface, err := cli.GetAddressUTXO(fromAddr, "")
		if err != nil {
			resp := ResFailed
			resp.Message = fmt.Sprintf("[%s] Get address [%s] utxo error: %+v", funcName, fromAddr, err)
			res.Status = resp
			return res
		}
		unspendUTXOList := utxoInterface.([]*client.UnspendUTXOList)
		if len(unspendUTXOList) < 1 {
			resp := ResFailed
			resp.Message = fmt.Sprintf("[%s] This address [%s] not have utxo", funcName, fromAddr)
			res.Status = resp
			return res
		}
		for _, unUtxo := range unspendUTXOList {
			if unUtxo.TxHash != vins[i].TxHash || unUtxo.Vout != vins[i].Vout {
				continue
			}
			useUTXOList = append(useUTXOList, unUtxo)
		}
	}

	// to地址数据处理
	var toAddrList []*client.ToAddrDetail
	for i, toAddr := range toAddrs {
		detail := client.GetToAddrDetail(toAddr, amounts[i])
		toAddrList = append(toAddrList, detail)
	}
	// 创建交易结构
	TxInfo, err := cli.BuildTransferInfoByList(useUTXOList, toAddrList, gasPrice, changeAddr)
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

/**
 * 合并签名并广播交易
 *
 * Params (chainName, priKeys []string, apiTx interface{})
 * chainName:
 *   链名称
 * priKeys:
 *   私钥
 * apiTx:
 *   交易信息
 */
func SignListAndSendTransferInfo(chainName string, priKeys []string, apiTx string) *CommonResp {
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

	// 签名并广播交易
	txHash, err := cli.SignListAndSendTransfer(apiTx, priKeys)
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
 * 签名数据
 *
 * Params (chainName, priKey string, apiTx interface{})
 * chainName:
 *   链名称
 * priKey:
 *   私钥
 * apiTx:
 *   交易信息
 */
func SignTransferInfo(chainName, priKey, apiTx string) *CommonResp {
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

	// 签名交易
	signRes, err := cli.SignTransferToRaw(apiTx, priKey)
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] sign transfer error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// 返回结果
	res.Status = ResSuccess
	res.Data = signRes
	return res
}

func MultiToMultiTransfer(chainName string, vins []ChooseUTXO, inputs []int64, toAddrs, amounts []string, gasPrice, changeAddr string) *CommonResp {
	res := &CommonResp{}
	funcName := "BuildTransferInfoByBTCList"

	// 链接节点
	cli, err := NewNodeService(chainName)
	defer cli.Close()
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] new node client error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// 最终使用UTXO
	var useUTXOList []*client.UnspendUTXOList
	for i, vin := range vins {
		rawAmount := decimal.NewFromInt(inputs[i]).BigInt()
		unUtxo := &client.UnspendUTXOList{
			TxHash:    vin.TxHash,
			Vout:      vin.Vout,
			RawAmount: rawAmount,
		}
		useUTXOList = append(useUTXOList, unUtxo)
	}

	// to地址数据处理
	var toAddrList []*client.ToAddrDetail
	for i, toAddr := range toAddrs {
		detail := client.GetToAddrDetail(toAddr, amounts[i])
		toAddrList = append(toAddrList, detail)
	}
	// 创建交易结构
	TxInfo, err := cli.BuildTransferInfoByList(useUTXOList, toAddrList, gasPrice, changeAddr)
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
