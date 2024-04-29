package main

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"syscall/js"
	"wallet_sdk"
	"wallet_sdk/client"
)

var (
	ErrorString = `{"code":1,"message":"%s"}`
)

func main() {
	InitJsFunc()
}

func InitJsFunc() {

	done := make(chan struct{}, 0)
	// 查询方法
	js.Global().Set("generateMnemonic", js.FuncOf(generateMnemonic))
	js.Global().Set("generateAccountByMnemonic", js.FuncOf(generateAccountByMnemonic))
	js.Global().Set("getBalance", js.FuncOf(getBalance))
	js.Global().Set("getTokenBalance", js.FuncOf(getTokenBalance))
	js.Global().Set("getTransferInfo", js.FuncOf(getTransferInfo))
	js.Global().Set("getContractInfo", js.FuncOf(getContractInfo))
	js.Global().Set("getGasPrice", js.FuncOf(getGasPrice))
	js.Global().Set("getNonce", js.FuncOf(getNonce))
	js.Global().Set("getBlockHeight", js.FuncOf(getBlockHeight))
	// 操作方法
	js.Global().Set("buildTransferInfo", js.FuncOf(buildTransferInfo))
	js.Global().Set("buildContractInfo", js.FuncOf(buildContractInfo))
	js.Global().Set("signAndSendTransferInfo", js.FuncOf(signAndSendTransferInfo))
	// 通用方法
	js.Global().Set("ethToGwei", js.FuncOf(ethToGwei))
	js.Global().Set("gweiToEth", js.FuncOf(gweiToEth))
	// 设置节点
	js.Global().Set("setNodeInfo", js.FuncOf(setNodeInfo))
	<-done
}

func returnResponse(res interface{}) string {
	data, err := json.Marshal(res)
	if err != nil {
		return fmt.Sprintf(ErrorString, err.Error())
	}
	return string(data)
}

func generateMnemonic(this js.Value, args []js.Value) interface{} {
	// 处理参数
	length := args[0].Int()
	language := args[1].String()
	res := wallet_sdk.GenerateMnemonic(length, language)
	return returnResponse(res)
}

func generateAccountByMnemonic(this js.Value, args []js.Value) interface{} {
	// 处理参数
	mnemonic := args[0].String()
	symbol := args[1].String()
	purpose := args[2].Int()
	var pp uint32
	if purpose != 0 {
		pp = uint32(purpose)
	}

	res := wallet_sdk.GenerateAccountByMnemonic(mnemonic, symbol, &pp)
	return returnResponse(res)
}

func getBalance(this js.Value, args []js.Value) interface{} {
	// 处理参数
	chain := args[0].String()
	address := args[1].String()
	// 异步方法需要回调函数
	funcName := args[2]

	go func() {
		defer recoverErr()
		res := wallet_sdk.GetBalanceByAddress(chain, address)
		// 使用回调函数处理结果
		funcName.Invoke(returnResponse(res))
	}()
	return nil
}

func getTokenBalance(this js.Value, args []js.Value) interface{} {
	// 处理参数
	chain := args[0].String()
	address := args[1].String()
	contract := args[2].String()
	// 异步方法需要回调函数
	funcName := args[3]

	go func() {
		defer recoverErr()
		res := wallet_sdk.GetBalanceByAddressAndContract(chain, address, contract)
		funcName.Invoke(returnResponse(res))
	}()
	return nil
}

func getTransferInfo(this js.Value, args []js.Value) interface{} {
	// 处理参数
	chain := args[0].String()
	txHash := args[1].String()
	// 异步方法需要回调函数
	funcName := args[2]

	go func() {
		defer recoverErr()
		res := wallet_sdk.GetTransaction(chain, txHash)
		funcName.Invoke(returnResponse(res))
	}()
	return nil
}

func getContractInfo(this js.Value, args []js.Value) interface{} {
	// 处理参数
	chain := args[0].String()
	contract := args[1].String()
	// 异步方法需要回调函数
	funcName := args[2]

	go func() {
		defer recoverErr()
		res := wallet_sdk.GetContractInfo(chain, contract)
		funcName.Invoke(returnResponse(res))
	}()
	return nil
}

func getGasPrice(this js.Value, args []js.Value) interface{} {
	// 处理参数
	chain := args[0].String()
	// 异步方法需要回调函数
	funcName := args[1]

	go func() {
		defer recoverErr()
		res := wallet_sdk.GetGasPrice(chain)
		funcName.Invoke(returnResponse(res))
	}()
	return nil
}

func getNonce(this js.Value, args []js.Value) interface{} {
	// 处理参数
	chain := args[0].String()
	address := args[1].String()
	params := args[2].String()
	// 异步方法需要回调函数
	funcName := args[3]

	go func() {
		defer recoverErr()
		res := wallet_sdk.GetNonce(chain, address, params)
		// 使用回调函数处理结果
		funcName.Invoke(returnResponse(res))
	}()
	return nil
}

func getBlockHeight(this js.Value, args []js.Value) interface{} {
	// 处理参数
	chain := args[0].String()
	// 异步方法需要回调函数
	funcName := args[1]

	go func() {
		defer recoverErr()
		res := wallet_sdk.GetBlockHeight(chain)
		// 使用回调函数处理结果
		funcName.Invoke(returnResponse(res))
	}()
	return nil
}

func buildTransferInfo(this js.Value, args []js.Value) interface{} {
	// 异步方法需要回调函数
	chain := args[0].String()
	fromAddr := args[1].String()
	toAddr := args[2].String()
	contract := args[3].String()
	amount := args[4].String()
	gasPrice := args[5].String()
	nonce := args[6].String()
	funcName := args[7]

	// 处理参数
	var list []interface{}
	for _, v := range args[8:] {
		list = append(list, v.String())
	}
	go func() {
		defer recoverErr()
		res := wallet_sdk.BuildTransferInfo(chain, fromAddr, toAddr, contract, amount, gasPrice, nonce)
		funcName.Invoke(returnResponse(res))
	}()
	return nil
}

func buildContractInfo(this js.Value, args []js.Value) interface{} {
	// 异步方法需要回调函数
	chain := args[0].String()
	contract := args[1].String()
	abiContent := args[2].String()
	gasPrice := args[3].String()
	nonce := args[4].String()
	params := args[5].String()
	funcName := args[6]

	// 处理参数
	var list []interface{}
	for _, v := range args[7:] {
		list = append(list, v.String())
	}
	go func() {
		defer recoverErr()
		res := wallet_sdk.BuildContractInfo(chain, contract, abiContent, gasPrice, nonce, params, list...)
		funcName.Invoke(returnResponse(res))
	}()
	return nil
}

func signAndSendTransferInfo(this js.Value, args []js.Value) interface{} {
	// 处理参数
	chain := args[0].String()
	priKey := args[1].String()
	txObj := args[2].String()
	// 异步方法需要回调函数
	funcName := args[3]

	go func() {
		defer recoverErr()
		res := wallet_sdk.SignAndSendTransferInfo(chain, priKey, txObj)
		// 使用回调函数处理结果
		funcName.Invoke(returnResponse(res))
	}()
	return nil
}

/* ------------------  常用方法  ---------------------------------------- */
func ethToGwei(this js.Value, args []js.Value) interface{} {
	// 处理参数
	value := args[0].String()
	return client.EthToGwei(value)
}
func gweiToEth(this js.Value, args []js.Value) interface{} {
	// 处理参数
	value := args[0].String()
	return client.WeiToGwei(value)
}
func setNodeInfo(this js.Value, args []js.Value) interface{} {
	// 处理参数
	node := client.Node{
		ChainType: wallet_sdk.ChainRelationForETH,
		Ip:        "http://192.168.10.173:8545",
		ChainId:   "11155111",
	}
	chainName := args[0].String()
	chainType := args[1].String()
	rpcURL := args[2].String()
	chainId := args[3].String()
	wallet_sdk.SetNodeInfo(chainName, node.ChainType, node.Ip, "", "", "", node.ChainId, "")
	return client.WeiToGwei(value)
}

func recoverErr() {
	if r := recover(); r != nil {
		//打印错误堆栈信息
		fmt.Printf("painc error: %+v\n", r)
		debug.PrintStack()
	}
}
