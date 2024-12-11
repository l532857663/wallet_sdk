package client

import (
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"math/big"
)

// 查询地址余额

func (c *BtcClient) GetBalance(addr, state string) (*big.Int, error) {
	utxoList := c.getAddressUTXO(addr)
	balance := big.NewInt(0)
	for _, utxo := range utxoList {
		balance.Add(balance, utxo.RawAmount)
	}
	return balance, nil
}

// 查询地址UTXO列表

func (c *BtcClient) GetAddressUTXO(addr, state string) (interface{}, error) {
	return c.getAddressUTXO(addr), nil
}

func (c *BtcClient) getAddressUTXO(address string) []*UnspendUTXOList {
	var res []*UnspendUTXOList
	addr, err := btcutil.DecodeAddress(address, c.Params)
	if err != nil {
		fmt.Printf("invalid recipet address: %v", err)
		return nil
	}
	fmt.Printf("wch----- addr: %+v\n", addr)
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
//func (c *BtcClient) getAddressUTXO(address string) []*UnspendUTXOList {
//	var res []*UnspendUTXOList
//	// 使用外部服务
//	addr, err := btcutil.DecodeAddress(address, c.Params)
//	if err != nil {
//		fmt.Printf("invalid recipet address: %v", err)
//		return nil
//	}
//	// 查询未花费的UTXO列表
//	unspendList, err := c.MempoolClient.ListUnspent(addr)
//	if err != nil {
//		fmt.Printf("GetListUnspent error: %+v", err)
//		return nil
//	}
//	if len(unspendList) == 0 {
//		fmt.Printf("no utxo for %v", addr)
//		return nil
//	}
//	// ScriptPubKey
//	spk, err := txscript.PayToAddrScript(addr)
//	if err != nil {
//		fmt.Printf("PayToAddrScript err: %v", err)
//		return nil
//	}
//	// 格式化
//	for _, unspend := range unspendList {
//		amount := unspend.Output.Value
//		tmp := &UnspendUTXOList{
//			TxHash:       unspend.Outpoint.Hash.String(),
//			ScriptPubKey: hex.EncodeToString(spk),
//			Vout:         unspend.Outpoint.Index,
//			Amount:       amount,
//			RawAmount:    big.NewInt(amount),
//		}
//		res = append(res, tmp)
//	}
//	return res
//}
