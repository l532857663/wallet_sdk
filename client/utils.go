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
	return wif.PrivKey, nil
}

func GetPrikeyByHex(prikey string) (*btcec.PrivateKey, error) {
	privateKeyBytes, err := hex.DecodeString(prikey)
	if err != nil {
		return nil, fmt.Errorf("GetPrikeyByHex DecodeString fatal, %+v", err)
	}
	privateKey, _ := btcec.PrivKeyFromBytes(privateKeyBytes)
	return privateKey, nil
}

func GetBTCAddress(key *btcec.PrivateKey, net *chaincfg.Params, addrType string) (string, error) {
	publicKey := key.PubKey().SerializeCompressed()
	pkHash := btcutil.Hash160(publicKey)
	switch addrType {
	case BTCAddrLegacy:
		addr, err := btcutil.NewAddressPubKeyHash(pkHash, net)
		return addr.EncodeAddress(), err
	case BTCAddrP2SH:
		addrP2WPKH, err := btcutil.NewAddressWitnessPubKeyHash(pkHash, net)
		script, err := txscript.PayToAddrScript(addrP2WPKH)
		addr, err := btcutil.NewAddressScriptHash(script, net)
		return addr.EncodeAddress(), err
	case BTCAddrP2WPKH:
		addr, err := btcutil.NewAddressWitnessPubKeyHash(pkHash, net)
		return addr.EncodeAddress(), err
	case BTCAddrP2TR:
		internalKey, err := btcec.ParsePubKey(publicKey)
		if err != nil {
			return "", nil
		}
		trpubKey := schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(internalKey))
		addr, err := btcutil.NewAddressTaproot(trpubKey, net)
		return addr.EncodeAddress(), err
	}
	return "", fmt.Errorf("Not get btc address")

	// tapHash, _ := hex.DecodeString("2c80c64d9f65be8794dbb43579129320051fdc3193a1ae9107202df3812aea57")
	// trwspubKey := schnorr.SerializePubKey(txscript.ComputeTaprootOutputKey(key.PubKey(), tapHash[:]))
	// addr3, _ := btcutil.NewAddressTaproot(trwspubKey, net)
	// fmt.Printf("BTC P2TP with script: %s\n", addr3.EncodeAddress())
	// tapPkScript, _ := txscript.PayToAddrScript(addr3)
	// fmt.Printf("the tappk script: %x\n", tapPkScript)
}

func CheckPrivateKey(prikey string) (*btcec.PrivateKey, error) {
	var privateKey *btcec.PrivateKey
	var err error
	keyLen := len(prikey)
	// 判断私钥是否是ETH系列的
	if keyLen == 66 && prikey[:2] == "0x" {
		privateKey, err = GetPrikeyByHex(prikey[2:])
	} else if keyLen == 64 {
		privateKey, err = GetPrikeyByHex(prikey)
	} else {
		wif, err := btcutil.DecodeWIF(prikey)
		if err != nil {
			return nil, err
		}
		privateKey = wif.PrivKey
	}
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}
