package client

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
)

// 签名并广播Taproot交易
func (c *BtcClient) SignAndSendTaprootTransfer(txObj, hexPrivateKey string, chainId *big.Int, idx int) (string, error) {
	txInfo := &BtcTransferInfo{}
	err := json.Unmarshal([]byte(txObj), txInfo)
	if err != nil {
		return "", err
	}
	apiTx := txInfo.ApiTx
	fmt.Printf("wch----- apiTx: %+v\n", apiTx)
	for idx, rti := range txInfo.UTXOList {
		prevOutScript, err := hex.DecodeString(rti.ScriptPubKey)
		if err != nil {
			fmt.Printf("invalid script key error: %+v\n", err)
			return "", err
		}
		fmt.Printf("wch------ prevOutScript: %+v\n", rti.ScriptPubKey)
		_, err = c.signTaprootUTXO(apiTx, hexPrivateKey, idx, prevOutScript)
		if err != nil {
			fmt.Printf("Sign err: %+v\n", err)
			return "", err
		}
	}
	raw, _ := getTxHex(apiTx)
	fmt.Printf("wch------ apiTx: %s, %s\n", apiTx.TxHash(), raw)
	return "", nil
}

func (c *BtcClient) signTaprootUTXO(apiTx *wire.MsgTx, privateKey string, utxoIdx int, utxoSciptPubKey []byte) (string, error) {
	txIn := apiTx.TxIn[utxoIdx]
	if nil == txIn {
		fmt.Printf("btc sign TxIn err! TxIn: %+v utxoIdx: %+v\n", apiTx.TxIn, utxoIdx)
		return "", fmt.Errorf("btc SignTx txIn is nil")
	}
	fmt.Printf("wch---- txIn: %+v\n", txIn)
	// 解析私钥
	wif, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "", fmt.Errorf("SignTx DecodeWIF fatal, " + err.Error())
	}
	if !wif.IsForNet(c.Params) {
		return "", fmt.Errorf("SignTx IsForNet not matched")
	}
	// hashCache := txscript.NewTxSigHashes(apiTx)
	// sig, err := txscript.SignTaprootOutput(&tx, hashCache, 0, int64(outputAmount), pkScript, txscript.SigHashAll, privKey, true)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// outPoint := txIn.PreviousOutPoint
	// prevOutputFetcher := txscript.NewMultiPrevOutFetcher(nil)
	// prevOutputFetcher.AddPrevOut(outPoint, apiTx.TxOut[utxoIdx])
	// txOut := prevOutputFetcher.FetchPrevOutput(txIn.PreviousOutPoint)
	// fmt.Printf("wch---- txOut: %+v\n", txOut)
	// witness, err := txscript.TaprootWitnessSignature(apiTx, txscript.NewTxSigHashes(apiTx, prevOutputFetcher),
	// 	utxoIdx, txOut.Value, txOut.PkScript, txscript.SigHashDefault, wif.PrivKey)
	// fmt.Printf("wch---- witness: %+v\n", witness)
	// apiTx.TxIn[utxoIdx].Witness = witness
	return "", nil
}
