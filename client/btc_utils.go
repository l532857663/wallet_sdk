package client

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"math/big"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// BTC的交易结构体

type BtcTransferInfo struct {
	ApiTx    *wire.MsgTx
	UTXOList []*UnspendUTXOList
}

func GetToAddrDetail(toAddr, amount string) *ToAddrDetail {
	rawAmount := BtcToSatoshi(amount)
	detail := &ToAddrDetail{
		Address:   toAddr,
		Amount:    amount,
		RawAmount: rawAmount,
	}
	return detail
}

// newPubkeyHash 生成公钥哈希脚本

func NewPubKeyHash(encodedAddr string, net *chaincfg.Params) ([]byte, error) {
	// Decode the recipent address.
	addr, err := btcutil.DecodeAddress(encodedAddr, net)
	if err != nil {
		return nil, fmt.Errorf("invalid recipet address: %w", err)
	}
	if !addr.IsForNet(net) {
		return nil, fmt.Errorf("%s is for the wrong network", encodedAddr)
	}
	// Create a new script which pays to the provided address.
	return txscript.PayToAddrScript(addr)
}

// PreCalculateSerializeSize
/* @Description: 预计交易的大小
 * @param apiTx *wire.MsgTx
 * @return int64
 */
func PreCalculateSerializeSize(apiTx *wire.MsgTx) int {
	return apiTx.SerializeSize() + BtcMaxSigScriptByteSize*len(apiTx.TxIn)
	//手续费计算规则： 148*len(apiTx.TxIn) + 34*len(apiTx.TxOut) + 10
}

func PayToPubKeyHashScript(pubKeyHash []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().AddOp(txscript.OP_DUP).AddOp(txscript.OP_HASH160).
		AddData(pubKeyHash).AddOp(txscript.OP_EQUALVERIFY).AddOp(txscript.OP_CHECKSIG).
		Script()
}

func PayToWitnessPubKeyHashScript(pubKeyHash []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(pubKeyHash).Script()
}

// Demo
func (c *BtcClient) getAddressUTXOForDemo(address string) []*UnspendUTXOList {
	var res []*UnspendUTXOList
	tmp := &UnspendUTXOList{
		TxHash:       "9a24b895b7bd528f1e24a7382792067272b6e8faa194ca35319d086604fc7fb6",
		ScriptPubKey: "76a914d8c9cf87df6269a9962023a57f18b93d1e4417fa88ac",
		Vout:         1,
		Amount:       7648,
		RawAmount:    big.NewInt(7648),
	}
	tmp1 := &UnspendUTXOList{
		TxHash:       "3c3c3bc3374297929eb5cd1c70e3d0525f6bf268f15c4c4855b14628b12dbe44",
		ScriptPubKey: "76a914d8c9cf87df6269a9962023a57f18b93d1e4417fa88ac",
		Vout:         1,
		Amount:       12080,
		RawAmount:    big.NewInt(12080),
	}
	res = append(res, tmp)
	res = append(res, tmp1)
	return res
}

func getTxHex(tx *wire.MsgTx) (string, error) {
	var buf bytes.Buffer
	if err := tx.Serialize(&buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

func GetAddressByPrivateKey(priKey *btcec.PrivateKey, params interface{}) (*btcutil.AddressPubKeyHash, error) {
	pubKey := priKey.PubKey()
	pkHash := btcutil.Hash160(pubKey.SerializeCompressed())
	addr, err := btcutil.NewAddressPubKeyHash(pkHash, params.(*chaincfg.Params))
	if err != nil {
		return nil, err
	}
	return addr, nil
}

func GetAddressByPKScript(pkScript []byte, params interface{}) (string, error) {
	_, addr, required, err := txscript.ExtractPkScriptAddrs(pkScript, params.(*chaincfg.Params))
	if err != nil {
		return "", err
	}
	if len(addr) == 0 || required == 0 {
		return "", fmt.Errorf("Not have address")
	}
	return addr[0].EncodeAddress(), nil
}
