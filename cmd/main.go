package main

import (
	"fmt"
	"wallet_sdk"
	"wallet_sdk/client"
)

func main() {
	// 生成助记词、使用助记词生成私钥、地址、生成seed
	testFunc()
	// 查询数据(自定义节点信息、余额、nonce、gas、contract、交易详情、块高等)
	// test1Func()
	// 导入私钥生成钱包信息
	// test2Func()
	// ETH的交易
	// test3Func()
	// BTC的交易
	// test4Func()
	// 部分签名出账
	// test5Func()
	// 多地址签名出账
	// test6Func()
}

var (
	length     = 12
	mnemonic   = "chaos coin cart couch system grunt soap never engine step glass quality"
	noContract = ""
	// ETH、BSC、POLYGON
	// addr   = "0x5827196b31CC0ddB815A7b297554916a76B6533A"
	// toAddr = "0x39447c3040124057147512c3D1477dAc339fcf8C"
	// toAddr = "0xd538657b3bd82d6ed51004e7e099b857384c4310"
	// toAddr = "0x77a50402d4d62a1b65f14ce79e8da0de9337d982"
	// toAddr = "0x8db4f1383517af7ae409d70f27495becfb3d45ee"
	// priKey    = "" //
	// chainName = wallet_sdk.ETH_Sepolia
	// chainName = wallet_sdk.BSC_Testnet
	// chainName = wallet_sdk.POLYGON_Testnet
	// contract = "0x4646f6a4c16788321bb1db1f904353c44f53fe1a" // ETH USDT
	// contract = "0x779877A7B0D9E8603169DdbD7836e478b4624789" // ETH LINK
	// contract = "0xb3066930566bEbdc96C2a94EC4aF9F6815Ec8004" // ETH SPC
	// contract  = "0xda977ea49bd752c4e2431ae41b31b1ca18c1226a" // BSC USDT
	// contract = "0xE209d8b4f1EF11802CFc76Bb4cD3a1C4762b6d37" // BSC USDT
	// contract = "0x2d7882bedcbfddce29ba99965dd3cdf7fcb10a1e" // POLYGON TST
	// contract = "0xd6afb55bcaa7711b7e6c0ea372a9156ce2d901de" // ETH TST
	// txHash = "0xa209c873736c117dde46418d0aa193a4c6d2207ee9fae458fe924aec0655679f" // MATIC
	// txHash = "0x847fa04f575e5cb47a0cfc413feb43975154c0d303ba117006dbb55aaec92646"

	// TRX
	// chainName = wallet_sdk.TRX_Nile
	// addr      = "TRtybNrManmeHwKHdehnfVumEjzXnsmrb3"
	// priKey    = ""
	// toAddr    = "TDUya7MQDTifg2EDyKZuCqScrnu2npnuor"
	// contract  = "TXLAQ63Xg1NAzckPwKHvzw7CSEmLMEqcdj" // TRX USDT

	// BTC
	chainName = wallet_sdk.BTC_Testnet
	// chainName = wallet_sdk.BTC_Testnet
	// chainName = wallet_sdk.BTC_Mainnet
	addr = "tb1pg0uc7ujx6rplw4wj73etg505jh49k63s7wc3kyngf73ze7ffue4skru6ld"
	// addr   = "2NBeoUKGLyk5ZfSDtAvsfWteYQaAKdUAniF"
	// addr   = "tb1pnkfdsmf6q4rmjtn3dunrmekxy57eq6xx7mnhthwu9z23u5hceluqvzmvul"
	toAddr    = "tb1qluw5qatmke9tqlu5wznp8u0r05lccpzhz23f5e"
	priKey    = ""
	priKeyHex = ""
	txHash    = "07b6936e51f8e7a83cd3f9fd811861224246924dd5b57052917aa87604fb2fa9"
	// addr      = "36F7BBBLxASGaAmgPnN15qLMUTnH7CTp16"
	// addr = "bc1ptz34pme4qp43qv6ykp3r0tqz4scn8frzg9e53m034w9st9ncpums67r7sv"
)

func testFunc() {
	// // 生成助记词
	// res := wallet_sdk.GenerateMnemonic(length, "")
	// fmt.Printf("res: %+v\n", res)
	// mnemonic = res.Data

	// 使用助记词生成账户地址、私钥
	var addressIndex uint32 = 0
	res1 := wallet_sdk.GenerateAccountByMnemonic(mnemonic, "BTCRegt", &addressIndex)
	fmt.Printf("res: %+v\n", res1)
	fmt.Printf("res Data: %+v\n", res1.Data)
	// 使用助记词生成seed
	// wallet_sdk.TestAccount(mnemonic, "BTC")

	return
}

func test1Func() {
	// 自定义节点信息
	chainName = "ethTest"
	node := client.Node{
		ChainType: wallet_sdk.MainCoinETH,
		Ip:        "http://192.168.10.173:8545",
		ChainId:   "11155111",
	}
	wallet_sdk.SetNodeInfo(chainName, node.ChainType, node.Ip, "", "", "", node.ChainId, "")
	// // 查询主币余额
	// res2 := wallet_sdk.GetBalanceByAddress(chainName, addr)
	// fmt.Printf("res: %+v\n", res2)
	// // // 查询代币余额
	// // res2 := wallet_sdk.GetBalanceByAddressAndContract(chainName, addr, contract)
	// // fmt.Printf("res: %+v\n", res2)
	// // 查询地址nonce
	// nonceData := wallet_sdk.GetNonce(chainName, addr, "latest")
	// fmt.Printf("res: %+v\n", nonceData)
	// // 查询节点gas price
	// gasPriceData := wallet_sdk.GetGasPrice(chainName)
	// fmt.Printf("res: %+v\n", gasPriceData.Data)
	// // 查询合约信息
	// res3 := wallet_sdk.GetContractInfo(chainName, contract)
	// fmt.Printf("res: %+v\n", res3)
	// fmt.Printf("res Data: %+v\n", res3.Data)

	// 查询交易详情
	res4 := wallet_sdk.GetTransaction(chainName, txHash)
	fmt.Printf("res: %+v\n", res4)
	fmt.Printf("res Data: %+v\n", res4.Data)
	fmt.Printf("res Data txInfo: %+v\n", res4.Data.TxInfo)

	// // 查询块高
	// res8 := wallet_sdk.GetBlockHeight(chainName)
	// fmt.Printf("res: %+v\n", res8)
}

func test2Func() {
	// 导入钱包操作
	res := wallet_sdk.ImportAddressByPrikey(priKey, "TRON")
	fmt.Printf("res: %+v\n", res)
	fmt.Printf("res Data: %+v\n", res.Data)
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

	// // 构建合约调用
	// res6 := wallet_sdk.BuildContractInfo(chainName, contract, abiContent, gasPrice, nonce, "approve", args...)
	// fmt.Printf("res: %+v\n", res6)

	// 签名交易
	signData := res5.Data
	res7 := wallet_sdk.SignTransferInfo(chainName, priKey, string(signData))
	// 签名并广播交易
	// res7 := wallet_sdk.SignAndSendTransferInfo(chainName, priKey, string(signData))
	fmt.Printf("res: %+v\n", res7)
}

func test4Func() {
	// 查询主币余额
	res2 := wallet_sdk.GetBalanceByAddress(chainName, addr)
	fmt.Printf("res: %+v\n", res2)
	// // 查询UTXO信息
	// addr = "tb1pfzl0rw44mkgevdauhrtzy5kdztjezyq0rnfqfppzxtnrwzdj553qvz6lux"
	// res2 := wallet_sdk.GetUTXOListByAddress(chainName, addr)
	// fmt.Printf("res: %+v\n", res2)

	// 查询节点gas price
	gasPriceData := wallet_sdk.GetGasPrice(chainName)
	// fmt.Printf("res: %+v\n", gasPriceData.Data)
	gasPrice := gasPriceData.Data.Average
	// gasPrice := "0.00000486"
	fmt.Printf("gasPrice: %+v\n", gasPrice)

	// 构建交易
	amount := "0.00000700"
	res5 := wallet_sdk.BuildTransferInfoByBTC(chainName, addr, toAddr, amount, gasPrice)
	fmt.Printf("res: %+v\n", res5)

	// 签名交易
	signData := res5.Data
	res7 := wallet_sdk.SignTransferInfo(chainName, priKey, string(signData))
	// 签名并广播交易
	// res7 := wallet_sdk.SignAndSendTransferInfo(chainName, priKey, string(signData))
	fmt.Printf("res: %+v\n", res7)
}

// 部分签名
func test5Func() {
	// // 查询地址余额
	// res2 := wallet_sdk.GetUTXOListByAddress(chainName, addr)
	// fmt.Printf("res: %+v\n", res2)
	// for _, utxo := range res2.Data.([]*client.UnspendUTXOList) {
	// 	fmt.Printf("utxo: %+v\n", utxo)
	// }
	// // 查询节点gas price
	gasPriceData := wallet_sdk.GetGasPrice(chainName)
	fmt.Printf("res: %+v\n", gasPriceData.Data)
	gasPrice := gasPriceData.Data.Average
	fmt.Printf("gasPrice: %+v\n", gasPrice)

	// 构建交易
	// gasPrice := "0.00000486"
	res5 := wallet_sdk.BuildPSBTransferInfo(chainName, priKey, gasPrice)
	fmt.Printf("res: %+v\n", res5)
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
