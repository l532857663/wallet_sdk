package main

import (
	"fmt"
	"wallet_sdk"
)

var (
	length   = 12
	mnemonic = "chaos coin cart couch system grunt soap never engine step glass quality"

	chainName, priKey, priKeyHex, addr, toAddr, contract, txHash string
)

func main() {
	// 生成助记词、使用助记词生成私钥、地址、生成seed
	//testMnemonicFunc()
	// 导入私钥生成钱包信息
	// testImportFunc()
	//testEthFunc()
	testBtcFunc()
	//testTronFunc()
	//testUtxoFunc()
}

func testUtxoFunc() {
	wallet_sdk.MustLoad("config.yml")
	wallet_sdk.InitNode()
	wallet_sdk.GetTransferByBlockHeight(1, 7150)
}
func testEthFunc() {
	// ETH、BSC、POLYGON
	// addr   = "0x5827196b31CC0ddB815A7b297554916a76B6533A"
	// toAddr = "0x39447c3040124057147512c3D1477dAc339fcf8C"
	// priKey    = ""
	//addr     = "0x3e7094B74549a6b8c4b4923cbC10Ef35c4D787Ce"
	contract = "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913" // BASE USDC
	// 查询数据(自定义节点信息、余额、nonce、gas、contract、交易详情、块高等)
	//test1Func()
	// 交易
	// test3Func()
	// 合约调用操作
	test7Func()
}
func testBtcFunc() {
	//chainName = wallet_sdk.BTC_Testnet
	chainName = wallet_sdk.BTC_RegTest
	addr = "2N9krkXFN1JZdSKSBKdYXCg9ZP8vUPvmo1e"
	toAddr = "bcrt1qyfre8aextm9dxj5vr45pfkycqft9qstkp8ufvj"
	// addr = "bc1ptz34pme4qp43qv6ykp3r0tqz4scn8frzg9e53m034w9st9ncpums67r7sv"
	txHash = "9f77dfe8c9709fd88d1f4ff1da022904fd8d36d07df68b285366b73a4c519405"
	priKey = "92K8ayCEhUEiv7JEGhxg9VSz2fQLEwRsk2ekJ5Eimsm1JbqKgjK"

	// 查询信息
	test2Func()
	// BTC的交易
	//test4Func()
	// 部分签名出账
	// test5Func()
	// 多地址签名出账
	// test6Func()
}

func testTronFunc() {
	// TRX
	// chainName = wallet_sdk.TRX_Nile
	// addr      = "TRtybNrManmeHwKHdehnfVumEjzXnsmrb3"
	// priKey    = ""
	// toAddr    = "TDUya7MQDTifg2EDyKZuCqScrnu2npnuor"
	// contract  = "TXLAQ63Xg1NAzckPwKHvzw7CSEmLMEqcdj" // TRX USDT
}
func testMnemonicFunc() {
	// // 生成助记词
	// res := wallet_sdk.GenerateMnemonic(length, "")
	// fmt.Printf("res: %+v\n", res)
	// mnemonic = res.Data

	// 使用助记词生成账户地址、私钥
	var addressIndex uint32 = 0
	res1 := wallet_sdk.GenerateAccountByMnemonic(mnemonic, "BTCRegT", &addressIndex)
	fmt.Printf("res: %+v\n", res1)
	fmt.Printf("res Data: %+v\n", res1.Data)
	// 使用助记词生成seed
	// wallet_sdk.TestAccount(mnemonic, "BTC")
	return
}

func testImportFunc() {
	// 导入钱包操作
	res := wallet_sdk.ImportAddressByPrikey(priKey, "TRON")
	fmt.Printf("res: %+v\n", res)
	fmt.Printf("res Data: %+v\n", res.Data)
}

func test1Func() {
	// 链接节点
	chainName = wallet_sdk.BSC_Testnet
	// 查询主币余额
	res2 := wallet_sdk.GetBalanceByAddress(chainName, addr)
	fmt.Printf("res balance: %+v\n", res2.Data)
	// 查询代币余额
	res2_1 := wallet_sdk.GetBalanceByAddressAndContract(chainName, addr, contract)
	fmt.Printf("res token: %+v\n", res2_1.Data)
	// 查询地址nonce
	nonceData := wallet_sdk.GetNonce(chainName, addr, "latest")
	fmt.Printf("res nonce: %+v\n", nonceData.Data)
	// 查询节点gas price
	gasPriceData := wallet_sdk.GetGasPrice(chainName)
	fmt.Printf("res gas: %+v\n", gasPriceData.Data)
	// 查询合约信息
	res3 := wallet_sdk.GetContractInfo(chainName, contract)
	// fmt.Printf("res: %+v\n", res3)
	fmt.Printf("res Data: %+v\n", res3.Data)

	// 查询交易详情
	res4 := wallet_sdk.GetTransaction(chainName, txHash)
	//fmt.Printf("res: %+v\n", res4)
	//fmt.Printf("res Data: %+v\n", res4.Data)
	fmt.Printf("res Data txInfo: %+v\n", res4.Data.TxInfo)

	// // 查询块高
	// res8 := wallet_sdk.GetBlockHeight(chainName)
	// fmt.Printf("res: %+v\n", res8)
}

func test2Func() {
	//// 查询块高
	//res1 := wallet_sdk.GetBlockHeight(chainName)
	//fmt.Printf("res: %+v\n", res1)
	//// 根据块高查询数据
	//height, _ := strconv.ParseInt(res1.Data, 10, 64)
	//res2 := wallet_sdk.GetBlockInfoByHeight(chainName, height)
	//fmt.Printf("res: %+v\n", res2)
	// 查询地址可用UTXO
	res3 := wallet_sdk.GetUTXOListByAddress(chainName, addr)
	fmt.Printf("res: %+v\n", res3)
	//// 查询交易详情
	//res4 := wallet_sdk.GetTransaction(chainName, txHash)
	////fmt.Printf("res: %+v\n", res4)
	////fmt.Printf("res Data: %+v\n", res4.Data)
	//fmt.Printf("res Data txInfo: %+v\n", res4.Data.TxInfo)
}

func test3Func() {
	// 查询地址nonce
	nonceData := wallet_sdk.GetNonce(chainName, addr, "latest")
	// fmt.Printf("res: %+v\n", nonceData)
	nonce := nonceData.Data
	// // 查询节点gas price
	// gasPriceData := wallet_sdk.GetGasPrice(chainName)
	// fmt.Printf("res: %+v\n", gasPriceData.Data)
	// gasPrice := gasPriceData.Data.High

	// nonce gasPrice
	gasPrice := "0.27"
	amount := "2000000000000000"

	// 构建交易
	res5 := wallet_sdk.BuildTransferInfo(chainName, addr, toAddr, amount, "", gasPrice, nonce)
	// res5 := wallet_sdk.BuildTransferInfo(chainName, addr, toAddr, amount, contract, gasPrice, nonce)
	fmt.Printf("res: %+v\n", res5)

	//// 构建合约调用
	//res6 := wallet_sdk.BuildContractInfo(chainName, contract, abiContent, gasPrice, nonce, "approve", args...)
	//fmt.Printf("res: %+v\n", res6)

	// 签名交易
	signData := res5.Data
	res7 := wallet_sdk.SignTransferInfo(chainName, priKey, string(signData))
	// 签名并广播交易
	// res7 := wallet_sdk.SignAndSendTransferInfo(chainName, priKey, string(signData))
	fmt.Printf("res: %+v\n", res7)
}

func test4Func() {
	//// 查询主币余额
	//res2 := wallet_sdk.GetBalanceByAddress(chainName, addr)
	//fmt.Printf("res: %+v\n", res2)
	// // 查询UTXO信息
	// addr = "tb1pfzl0rw44mkgevdauhrtzy5kdztjezyq0rnfqfppzxtnrwzdj553qvz6lux"
	// res2 := wallet_sdk.GetUTXOListByAddress(chainName, addr)
	// fmt.Printf("res: %+v\n", res2)

	//// 查询节点gas price
	//gasPriceData := wallet_sdk.GetGasPrice(chainName)
	//// fmt.Printf("res: %+v\n", gasPriceData.Data)
	//gasPrice := gasPriceData.Data.Average
	gasPrice := "0.00000486"
	fmt.Printf("gasPrice: %+v\n", gasPrice)

	// 构建交易
	amount := "0.00007000"
	res5 := wallet_sdk.BuildTransferInfoByBTC(chainName, addr, toAddr, amount, gasPrice)
	fmt.Printf("res: %+v\n", res5)

	// 签名交易
	signData := res5.Data
	res7 := wallet_sdk.SignTransferInfo(chainName, priKey, signData)
	// 签名并广播交易
	// res7 := wallet_sdk.SignAndSendTransferInfo(chainName, priKey, string(signData))
	fmt.Printf("res: %+v\n", res7)
}

// 部分签名
func test5Func() {
	//// 查询地址余额
	//res2 := wallet_sdk.GetUTXOListByAddress(chainName, addr)
	//fmt.Printf("res: %+v\n", res2)
	//for _, utxo := range res2.Data.([]*client.UnspendUTXOList) {
	//	fmt.Printf("utxo: %+v\n", utxo)
	//}
	// // 查询节点gas price
	gasPriceData := wallet_sdk.GetGasPrice(chainName)
	fmt.Printf("res: %+v\n", gasPriceData.Data)
	gasPrice := gasPriceData.Data.Average
	fmt.Printf("gasPrice: %+v\n", gasPrice)

	// 构建交易
	// gasPrice := "0.00000486"
	//res5 := wallet_sdk.BuildPSBTransferInfo(chainName, priKey, gasPrice)
	//fmt.Printf("res: %+v\n", res5)
}

func test6Func() {
	// 多地址签名出账
	// 查询节点gas price
	gasPriceData := wallet_sdk.GetGasPrice(chainName)
	// fmt.Printf("res: %+v\n", gasPriceData.Data)
	gasPrice := gasPriceData.Data.Average
	fmt.Printf("gasPrice: %+v\n", gasPrice)

	// 构建交易
	froms := []string{"n1HE1YJ1zF5U5aiX2DNu5WhjE9KFrkSKkx", "mqrg3rNg7cCLHVRCqYpzQoNE744DvJreeN"}
	utxos := []wallet_sdk.ChooseUTXO{
		{TxHash: "303b1e8f45adf1b1d07e1febb8fe0da2e4772862bf4189fbb120c188c5ecd95b", Vout: 1},
		{TxHash: "7fc39b92f2bc4bd12e5c441bf7e9f3f56cae02995c32ce07d0a14dc1b7ae872c", Vout: 0},
	}
	toAddrs := []string{"n4R9ztyWCfkuoX3vmWYwbYcohJRbL4yao1", "n4R9ztyWCfkuoX3vmWYwbYcohJRbL4yao1"}
	amounts := []string{"0.000003", "0.000003"}
	changeAddr := ""
	res5 := wallet_sdk.BuildTransferInfoByBTCList(chainName, froms, utxos, toAddrs, amounts, gasPrice, changeAddr)
	fmt.Printf("res: %+v\n", res5)

	// 签名并广播交易
	priKeys := []string{"cQSreoKBANpfNxLHD6v1crHE3rz44Q7hZPsV2XaJVQv6dA5eXGQV", "cTAroTKYviiVEqxPmTW43JEU56EFZLhLWYgcxHk8nfBovj72eXbT"}
	signData := res5.Data
	res7 := wallet_sdk.SignListAndSendTransferInfo(chainName, priKeys, string(signData))
	fmt.Printf("res: %+v\n", res7)
}

func test7Func() {
	// 链接节点
	chainName = wallet_sdk.BSC_Testnet
	addr = "0x39447c3040124057147512c3D1477dAc339fcf8C"

	// 查询主币余额
	res2 := wallet_sdk.GetBalanceByAddress(chainName, addr)
	fmt.Printf("res balance: %+v\n", res2.Data)

	//// 构建合约调用 BSC testnet v2 Router
	//amount := "100"
	//getAmountsOutPath := `["0xAFbcDd676D8BAD865B50Af6F77bD4914a15c3F70", "0xDA9dbD3c0C9e613973559f9593aB8d6cEd579CDD"]`
	//res6 := wallet_sdk.GetContractInfoByFunc(chainName, "0xD99D1c33F9fC3444f8101754aBC46c52416550D1", "getAmountsOut", amount, getAmountsOutPath)
	//fmt.Printf("res: %+v\n", res6)

	// 构建合约调用 BSC testnet v2 Factory
	tokenA := "0xAFbcDd676D8BAD865B50Af6F77bD4914a15c3F70"
	tokenB := "0xDA9dbD3c0C9e613973559f9593aB8d6cEd579CDD"
	res7 := wallet_sdk.GetContractInfoByFunc(chainName, "0x6725F303b657a9451d8BA641348b6761A6CC7a17", "getPair", tokenA, tokenB)
	fmt.Printf("res: %+v\n", res7)
}
