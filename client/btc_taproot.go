package client

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// 签名并广播Taproot交易
func (c *BtcClient) SignAndSendTaprootTransfer(txObj string, privateKey string, chainId *big.Int, idx int) (string, error) {
	txInfo := &BtcTransferInfo{}
	err := json.Unmarshal([]byte(txObj), txInfo)
	if err != nil {
		return "", err
	}
	// 解析私钥
	wif, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "", fmt.Errorf("SignTx DecodeWIF fatal, " + err.Error())
	}
	if !wif.IsForNet(c.Params) {
		return "", fmt.Errorf("SignTx IsForNet not matched")
	}
	apiTx := txInfo.ApiTx
	fmt.Printf("wch----- apiTx: %+v\n", apiTx)
	prevOutFetcher := txscript.NewMultiPrevOutFetcher(nil)
	var privateKeys []*btcec.PrivateKey
	for i := 0; i < len(apiTx.TxIn); i++ {
		input := apiTx.TxIn[i]
		utxoInfo := txInfo.UTXOList[i]
		outPoint := &input.PreviousOutPoint
		pkScript, err := hex.DecodeString(utxoInfo.ScriptPubKey)
		if err != nil {
			return "", err
		}
		txOut := wire.NewTxOut(utxoInfo.Amount.CoefficientInt64(), pkScript)
		prevOutFetcher.AddPrevOut(*outPoint, txOut)
		privateKeys = append(privateKeys, wif.PrivKey)
	}

	if err := Sign(apiTx, privateKeys, prevOutFetcher); err != nil {
		return "", err
	}

	raw, _ := getTxHex(apiTx)
	fmt.Printf("apiTx txHash: %s, info: %s\n", apiTx.TxHash(), raw)

	txHash, err := c.Client.SendRawTransaction(apiTx, false)
	if nil != err {
		return "", fmt.Errorf("Broadcast SendRawTransaction fatal, " + err.Error())
	}
	return txHash.String(), nil
}
