package wallet_sdk

import (
	"math/big"
	"wallet_sdk/client"
)

var (
	Ins = []client.Input{
		{
			// 1000
			TxId: "1d3b6d58e81c380301d77a024fa358eaefd8883f63e94ec205bea1ca16cd6cdc",
			VOut: 0,
		},
		{
			// 14000
			TxId: "303b1e8f45adf1b1d07e1febb8fe0da2e4772862bf4189fbb120c188c5ecd95b",
			VOut: 1,
		},
	}
	Outs = []client.Output{
		{
			Address: "tb1qtqdffxec6drzm7je7nkpdxkcum6nf3mtfqknrq",
			Amount:  3300000,
		},
		{
			Address: "n1HE1YJ1zF5U5aiX2DNu5WhjE9KFrkSKkx",
			Amount:  72688,
		},
	}
)

var (
	a1 = &client.UnspendUTXOList{
		TxHash:       "d39462f58bf84a88878edc79e1f0707d374b1d4ac8f4b392efd83d9482464d79",
		ScriptPubKey: "76a914d8c9cf87df6269a9962023a57f18b93d1e4417fa88ac",
		Vout:         1,
		Amount:       833,
		RawAmount:    big.NewInt(833),
	}
	a2 = &client.UnspendUTXOList{
		TxHash:       "8788861fe7966106313104a14e9b274881f8af561053b0ee5576d9bd3b0943a9",
		ScriptPubKey: "76a914d8c9cf87df6269a9962023a57f18b93d1e4417fa88ac",
		Vout:         1,
		Amount:       14333,
		RawAmount:    big.NewInt(14333),
	}
	b1 = &client.UnspendUTXOList{
		TxHash:       "1d3b6d58e81c380301d77a024fa358eaefd8883f63e94ec205bea1ca16cd6cdc",
		ScriptPubKey: "512043f98f7246d0c3f755d2f472b451f495ea5b6a30f3b11b12684fa22cf929e66b",
		Vout:         0,
		Amount:       1000,
		RawAmount:    big.NewInt(1000),
	}
	b2 = &client.UnspendUTXOList{
		TxHash:       "6ffa44f0cd70cdd611885419a0a67f63731e374fb591757d8824a7f7306fbe39",
		ScriptPubKey: "512043f98f7246d0c3f755d2f472b451f495ea5b6a30f3b11b12684fa22cf929e66b",
		Vout:         0,
		Amount:       546,
		RawAmount:    big.NewInt(546),
	}
)
