package wallet

import "github.com/btcsuite/btcd/btcec/v2"

func NewWalletByPrivateKey(privatekey *btcec.PrivateKey, opts ...Option) (Wallet, error) {
	var (
		o = newOptions(opts...)
	)
	key := &Key{
		Opt:          o,
		Private:      privatekey,
		Public:       privatekey.PubKey(),
		PrivateECDSA: privatekey.ToECDSA(),
		PublicECDSA:  &privatekey.ToECDSA().PublicKey,
	}

	coin, ok := coins[key.Opt.CoinType]
	if !ok {
		return nil, ErrCoinTypeUnknow
	}

	return coin(key), nil
}
