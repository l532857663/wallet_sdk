package client

import "C"
import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/shopspring/decimal"
	"log"
	"math/big"
	"wallet_sdk/utils/dir"
)

// 查询地址余额

func (c *BtcClient) GetBalance(addr, state string) (*big.Int, error) {
	utxoList := getAddressUTXO(c, addr)
	balance := decimal.Zero.BigInt()
	for _, utxo := range utxoList {
		fmt.Printf("utxo: %+v\n", utxo)
		balance.Add(balance, utxo.RawAmount)
	}
	return balance, nil
}

// 查询地址UTXO列表

func (c *BtcClient) GetAddressUTXO(addr, state string) (interface{}, error) {
	return getAddressUTXO(c, addr), nil
}

func getAddressUTXO(c *BtcClient, addr string) []*UnspendUTXOList {
	switch c.Params {
	case &chaincfg.MainNetParams, &chaincfg.TestNet3Params:
		return getMempoolAddressUTXO(c, addr)
	case &chaincfg.RegressionNetParams:
		list, err := dir.ReadFiles(c.UnSpendUtxoPath + "/" + addr)
		if err != nil {
			log.Printf("GetFiles error: %v", err)
			return nil
		}
		var res []*UnspendUTXOList
		for _, data := range list {
			fmt.Printf("data: %+v\n", string(data))
			tmp := &UnspendUTXOList{}
			json.Unmarshal(data, &tmp)
			fmt.Printf("tmp: %+v\n", tmp)
			tmp.RawAmount = tmp.Amount.BigInt()
			res = append(res, tmp)
		}
		return res
	}
	log.Fatalf("Params [%v] error", c.Params)
	return nil
}

// 节点钱包查询UTXO列表
func getLocalWalletAddressUTXO(c *BtcClient, address string) []*UnspendUTXOList {
	var res []*UnspendUTXOList
	addr, err := btcutil.DecodeAddress(address, c.Params)
	if err != nil {
		fmt.Printf("invalid recipet address: %v", err)
		return nil
	}
	fmt.Printf("addr: %+v\n", addr)
	addrList := []btcutil.Address{
		addr,
	}
	unSpentList, err := c.Client.ListUnspentMinMaxAddresses(1, 9999999, addrList)
	if err != nil {
		fmt.Printf("Get ListUnspentMinMaxAddresses error: %v", err)
	}
	for _, unSpent := range unSpentList {
		fmt.Printf("wch==== balances: %+v\n", unSpent)
	}
	return res
}

// 使用mempool查询可用UTXO
func getMempoolAddressUTXO(c *BtcClient, address string) []*UnspendUTXOList {
	var res []*UnspendUTXOList
	// 使用外部服务
	addr, err := btcutil.DecodeAddress(address, c.Params)
	if err != nil {
		fmt.Printf("invalid recipet address: %v", err)
		return nil
	}
	// 查询未花费的UTXO列表
	unspendList, err := c.MempoolClient.ListUnspent(addr)
	if err != nil {
		fmt.Printf("GetListUnspent error: %+v", err)
		return nil
	}
	if len(unspendList) == 0 {
		fmt.Printf("no utxo for %v", addr)
		return nil
	}
	// ScriptPubKey
	spk, err := txscript.PayToAddrScript(addr)
	if err != nil {
		fmt.Printf("PayToAddrScript err: %v", err)
		return nil
	}
	// 格式化
	for _, unspend := range unspendList {
		amount := decimal.NewFromInt(unspend.Output.Value)
		tmp := &UnspendUTXOList{
			TxHash:       unspend.Outpoint.Hash.String(),
			ScriptPubKey: hex.EncodeToString(spk),
			Vout:         unspend.Outpoint.Index,
			Amount:       amount,
			RawAmount:    amount.BigInt(),
		}
		res = append(res, tmp)
	}
	return res
}
