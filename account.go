package wallet_sdk

import (
	"encoding/base64"
	"fmt"
	"wallet_sdk/client"
	"wallet_sdk/utils"
	"wallet_sdk/wallet"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/crypto"
)

/**
 * 创建助记词
 *
 * Params (length int, langauge string)
 * length:
 *   生成助记词长度 使用 12、24
 * language:
 *   使用语言，默认英文 (简体中文：chinese_simplified、繁体中文：chinese_traditional)
 */

func GenerateMnemonic(length int, langauge string) *CommonResp {
	res := &CommonResp{}
	funcName := "GenerateMnemonic"

	// 生成助记词
	mnemonic, err := wallet.NewMnemonic(length, langauge)
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] new mnemonic error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// 返回结果
	res.Data = mnemonic

	return res
}

/**
 * 使用助记词生成账号
 *
 * Params (mnemonic, symbol string)
 * mnemonic:
 *   助记词字符串
 * symbol:
 *   选择的链类型 e.g.： "BTC"、"ETH"
 */

func GenerateAccountByMnemonic(mnemonic, symbol string) *AccountInfoResp {
	res := &AccountInfoResp{}
	funcName := "GenerateAccountByMnemonic"

	// 生成构建账户的结构体
	master, err := wallet.NewKey(
		wallet.Mnemonic(mnemonic),
	)
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] new key error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// 获取coinType
	coinType := utils.Symbol2CoinType(symbol)
	if coinType == wallet.Zero {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] this chain [%s] is not supported for now", funcName, symbol)
		res.Status = resp
		return res
	}

	// 生成账户信息
	var account wallet.Wallet
	if symbol == "BTCTest" {
		account, err = master.GetWallet(wallet.CoinType(coinType), wallet.Params(&wallet.BTCTestnetParams))
	} else {
		account, err = master.GetWallet(wallet.CoinType(coinType))
	}
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] get account error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// 账户地址
	address, err := account.GetAddress()
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] get address error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// 私钥
	priKey, err := account.GetPrivateKey()
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] get private key error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	accountInfo := &AccountInfo{
		Address:    address,
		PrivateKey: priKey,
	}

	// 返回结果
	res.Status = ResSuccess
	res.Data = accountInfo

	return res
}

/**
 * 导入私钥生成账号
 *
 * Params (prikey, symbol string)
 * prikey:
 *   私钥字符串
 * symbol:
 *   选择的链类型 e.g.： "BTC"、"ETH"
 */

func ImportAddressByPrikey(prikey, symbol string) *AccountInfoResp {
	res := &AccountInfoResp{}
	funcName := "ImportAddressByPrikey"

	if len(prikey) != 64 {
		if len(prikey) == 66 && prikey[:2] == "0x" {
			prikey = prikey[2:]
		} else {
			resp := ResFailed
			resp.Message = fmt.Sprintf("[%s] the private_key error", funcName)
			res.Status = resp
			return res
		}
	}

	// 私钥
	priKey, err := crypto.HexToECDSA(prikey)
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] get address error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// 账户地址
	pubkey := (*priKey).PublicKey
	address := crypto.PubkeyToAddress(pubkey).Hex()

	accountInfo := &AccountInfo{
		Address:    address,
		PrivateKey: prikey,
	}

	// 返回结果
	res.Status = ResSuccess
	res.Data = accountInfo

	return res
}

func TestAccount(mnemonic, symbol string) {
	funcName := "TestAccount"
	seed, err := wallet.NewSeed(mnemonic, "", "")
	fmt.Printf("wch-------- seed: %d, %s, %x\n", len(seed), base64.StdEncoding.EncodeToString(seed), seed)
	// 获取coinType
	coinType := utils.Symbol2CoinType(symbol)
	// 生成构建账户的结构体
	master, err := wallet.NewKey(
		wallet.Mnemonic(mnemonic),
	)
	if err != nil {
		fmt.Printf("[%s]NewKey error: %+v\n", funcName, err)
		return
	}

	// // 生成账户信息
	for i := 0; i <= 2; i++ {
		account, err := master.GetWallet(wallet.CoinType(coinType), wallet.AddressIndex(uint32(i)), wallet.Params(&wallet.BTCParams))
		if err != nil {
			fmt.Printf("[%s]GetWallet error: %+v\n", funcName, err)
			return
		}
		fmt.Printf("name: %+v\n", account.GetName())
		address, _ := account.GetAddress()
		fmt.Printf("wch----- address: %s\n", address)
		priStr, _ := account.GetPrivateKey()
		fmt.Printf("wch----- priStr: %s\n", priStr)
		priHex := account.GetKey().PrivateHex()
		fmt.Printf("wch----- priKey: %s\n", priHex)
		pubKey, _ := account.GetKey().Extended.Neuter()
		fmt.Printf("wch----- pubKey: %s\n", pubKey.String())
	}
}

func GetPrikeyAndPubkey(chainName, wifKey, hexKey string) {
	funcName := "GetPrikeyAndPubkeyByWIF"
	net := &chaincfg.TestNet3Params
	if wifKey != "" {
		key, _ := client.GetPrikeyByWIF(wifKey, net)
		client.GetBTCAddress(key, net)
	} else if hexKey != "" {
		key, _ := client.GetPrikeyByHex(hexKey, net)
		client.GetBTCAddress(key, net)
	} else {
		fmt.Printf("[%s] params error %s, %s!", funcName, wifKey, hexKey)
	}
	return
}
