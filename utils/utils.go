package utils

import (
	"math/big"
	"strings"
	"wallet_sdk/wallet"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
)

func Symbol2CoinType(symbol string) uint32 {
	switch symbol {
	case "BTC", "BTCTest", "BTCRegt":
		return wallet.BTC
	case "ETH":
		return wallet.ETH
	case "TRON":
		return wallet.TRX
	case "SOLANA":
		return wallet.SOL
	default:
		// 暂不支持改链
		return wallet.Zero
	}
}

func HexTobigInt(v string) *big.Int {
	bv := big.NewInt(0)
	if ok := strings.HasPrefix(v, "0x"); !ok {
		return nil
	}
	if len(v) < 3 {
		return nil
	}
	d, ok := bv.SetString(v[2:], 16)
	if !ok {
		return nil
	}
	return d
}

func ByteTobigInt(vByte []byte) *big.Int {
	v := hexutil.Encode(vByte)
	if len(v) < 3 {
		return big.NewInt(0)
	}
	return HexTobigInt(v)
}

func StringToDecimal(a string) (decimal.Decimal, error) {
	value, err := decimal.NewFromString(a)
	if err != nil {
		return decimal.NewFromInt(0), err
	}
	return value, nil
}

func StringTobigInt(a string) *big.Int {
	value, ok := new(big.Int).SetString(a, 10)
	if !ok {
		return big.NewInt(0)
	}
	return value
}

func Int64ToSatoshi(amount int64) btcutil.Amount {
	return btcutil.Amount(amount)
}
