package wallet_sdk

import (
	"encoding/base64"
	"fmt"
	"wallet_sdk/client"
	"wallet_sdk/utils"
	"wallet_sdk/wallet"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/chaincfg"
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

func GenerateAccountByMnemonic(mnemonic, symbol string, addressIndex *uint32) *AccountInfoResp {
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
	master.Opt.CoinType = coinType

	// 生成账户信息
	if addressIndex != nil {
		master.Opt.AddressIndex = *addressIndex
	}
	if symbol == "BTCTest" {
		master.Opt.Params = &wallet.BTCTestnetParams
	} else if symbol == "BTCRegt" {
		master.Opt.Params = &wallet.BTCRegtestParams
	}
	account, err := master.GetWallet(wallet.CoinType(master.Opt.CoinType), wallet.Params(master.Opt.Params), wallet.AddressIndex(master.Opt.AddressIndex))
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
	// BTC多地址处理
	if account.GetSymbol() == MainCoinBTC {
		accountInfo.BtcAddrList = getBtcMultiAddr(master)
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

	key, err := client.CheckPrivateKey(prikey)
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] check private key error: %+v", funcName, err)
		res.Status = resp
		return res
	}
	fmt.Printf("wch---- key: %+v\n", key)

	// 获取coinType
	coinType := utils.Symbol2CoinType(symbol)
	if coinType == wallet.Zero {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] this chain [%s] is not supported for now", funcName, symbol)
		res.Status = resp
		return res
	}

	account, err := wallet.NewWalletByPrivateKey(key, wallet.CoinType(coinType))
	if err != nil {
		resp := ResFailed
		resp.Message = fmt.Sprintf("[%s] new wallet by private key error: %+v", funcName, err)
		res.Status = resp
		return res
	}

	// 账户地址
	address, _ := account.GetAddress()
	// 私钥
	priKey, _ := account.GetPrivateKey()
	accountInfo := &AccountInfo{
		Address:    address,
		PrivateKey: priKey,
	}
	// 网络类型
	decoded := base58.Decode(prikey)
	if symbol == MainCoinBTC && len(decoded) > 0 {
		net, ok := wallet.GetParamsList[decoded[0]]
		if ok {
			wif, _ := btcutil.NewWIF(key, &net, false)
			for _, t := range client.BTCAddrList {
				addr, _ := client.GetBTCAddress(key, &net, t)
				tmp := BtcAddress{
					Address:     addr,
					AddressType: t,
					PrivateKey:  wif.String(),
				}
				accountInfo.BtcAddrList = append(accountInfo.BtcAddrList, tmp)
			}
		}
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

func GetPrikeyAndPubkey(chainName, wifKey, hexKey string) string {
	funcName := "GetPrikeyAndPubkeyByWIF"
	net := &chaincfg.MainNetParams
	if chainName == BTC_Testnet {
		net = &chaincfg.TestNet3Params
	}
	addr := ""
	if wifKey != "" {
		key, _ := client.GetPrikeyByWIF(wifKey, net)
		addr, _ = client.GetBTCAddress(key, net, client.BTCAddrLegacy)
	} else if hexKey != "" {
		key, _ := client.GetPrikeyByHex(hexKey)
		addr, _ = client.GetBTCAddress(key, net, client.BTCAddrLegacy)
	} else {
		fmt.Printf("[%s] params error %s, %s!", funcName, wifKey, hexKey)
	}
	return addr
}

func getBtcMultiAddr(master *wallet.Key) []BtcAddress {
	var btcAddrList []BtcAddress
	for i, t := range client.BTCAddrList {
		master.Opt.Purpose = wallet.BtcPurposeList[i]
		account, err := master.GetWallet(wallet.Purpose(master.Opt.Purpose), wallet.CoinType(master.Opt.CoinType), wallet.Params(master.Opt.Params), wallet.AddressIndex(master.Opt.AddressIndex))
		if err != nil {
			continue
		}
		key := account.GetKey().Private
		net := master.Opt.Params
		wif, _ := btcutil.NewWIF(key, net, false)
		addr, _ := client.GetBTCAddress(key, net, t)
		tmp := BtcAddress{
			Address:     addr,
			AddressType: t,
			PrivateKey:  wif.String(),
		}
		btcAddrList = append(btcAddrList, tmp)
	}
	return btcAddrList
}
