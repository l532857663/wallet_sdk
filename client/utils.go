package client

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"wallet_sdk/utils"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/shopspring/decimal"
)

var (
	paramsErr = func(name string) error { return fmt.Errorf("The [%+v] params error!", name) }
)

func EthToGwei(value string) string {
	v, err := utils.StringToDecimal(value)
	if err != nil {
		return err.Error()
	}
	return v.Mul(EthBaseDecimal).String()
}

func WeiToGwei(value string) string {
	v, err := utils.StringToDecimal(value)
	if err != nil {
		return err.Error()
	}
	return v.Div(EthBaseDecimal).String()
}

func BtcToSatoshi(value string) *big.Int {
	v, err := utils.StringToDecimal(value)
	if err != nil {
		return big.NewInt(0)
	}
	return v.Mul(BtcBaseDecimal).BigInt()
}

func SatoshiToBtc(value string) decimal.Decimal {
	v, err := utils.StringToDecimal(value)
	if err != nil {
		return v
	}
	return v.Div(BtcBaseDecimal)
}

type UnspendUTXOs []*UnspendUTXOList

func (u UnspendUTXOs) Len() int {
	return len(u)
}
func (u UnspendUTXOs) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}
func (u UnspendUTXOs) Less(i, j int) bool {
	return u[i].Amount < u[j].Amount
}

func DescSortUnspendUTXO(data []*UnspendUTXOList) {
	// 获取对应map和数量slice
	sort.Sort(UnspendUTXOs(data))
}

func GetPrikeyByWIF(prikey string, net *chaincfg.Params) (*btcec.PrivateKey, error) {
	wif, err := btcutil.DecodeWIF(prikey)
	if err != nil {
		return nil, fmt.Errorf("GetPrikeyByWIF DecodeWIF fatal, %+v", err)
	}
	if !wif.IsForNet(net) {
		return nil, fmt.Errorf("GetPrikeyByWIF IsForNet not matched")
	}
	fmt.Printf("Hex key: %+v\n", hex.EncodeToString(wif.PrivKey.Serialize()))
	fmt.Printf("pubKey compressed: %+v\n", hex.EncodeToString(wif.PrivKey.PubKey().SerializeCompressed()))
	fmt.Printf("pubKey uncompressed: %+v\n", hex.EncodeToString(wif.PrivKey.PubKey().SerializeUncompressed()))
	return wif.PrivKey, nil
}

func GetPrikeyByHex(prikey string, net *chaincfg.Params) (*btcec.PrivateKey, error) {
	privateKeyBytes, err := hex.DecodeString(prikey)
	if err != nil {
		return nil, fmt.Errorf("GetPrikeyByHex DecodeString fatal, %+v", err)
	}
	privateKey, publicKey := btcec.PrivKeyFromBytes(privateKeyBytes)
	wif, _ := btcutil.NewWIF(privateKey, net, true)
	fmt.Printf("WIF key: %+v\n", wif.String())
	fmt.Printf("pubKey compressed: %+v\n", hex.EncodeToString(publicKey.SerializeCompressed()))
	fmt.Printf("pubKey uncompressed: %+v\n", hex.EncodeToString(publicKey.SerializeUncompressed()))
	return wif.PrivKey, nil
}

func GetBTCAddress(key *btcec.PrivateKey, net *chaincfg.Params) {
	pkHash := btcutil.Hash160(key.PubKey().SerializeCompressed())
	addr, _ := btcutil.NewAddressPubKeyHash(pkHash, net)
	fmt.Printf("BTC P2PKH: %s\n", addr.EncodeAddress())

	// script, _ := txscript.PayToAddrScript(addr)
	// addrX, _ := btcutil.NewAddressScriptHash(script, net)
	// fmt.Printf("BTC P2SH: %s\n", addrX.EncodeAddress())

	addr1, _ := btcutil.NewAddressWitnessPubKeyHash(pkHash, net)
	fmt.Printf("BTC P2WPKH: %s\n", addr1.EncodeAddress())

	publicKey := key.PubKey().SerializeCompressed()
	internalKey, err := btcec.ParsePubKey(publicKey)
	if err != nil {
		return
	}
	trpubKey := schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(internalKey))
	addr2, _ := btcutil.NewAddressTaproot(trpubKey, net)
	fmt.Printf("BTC P2TP: %s\n", addr2.EncodeAddress())

	tapHash, _ := hex.DecodeString("2c80c64d9f65be8794dbb43579129320051fdc3193a1ae9107202df3812aea57")
	trwspubKey := schnorr.SerializePubKey(txscript.ComputeTaprootOutputKey(key.PubKey(), tapHash[:]))
	addr3, _ := btcutil.NewAddressTaproot(trwspubKey, net)
	fmt.Printf("BTC P2TP with script: %s\n", addr3.EncodeAddress())
	tapPkScript, _ := txscript.PayToAddrScript(addr3)
	fmt.Printf("the tappk script: %x\n", tapPkScript)
}
