package elastic

import "github.com/shopspring/decimal"

type ElasticConfig struct {
	Host     string `mapstructure:"host"     json:"host"     yaml:"host"`
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}

// *********************************************** BTC UTXO info ***********************************************
type AddressUTXOInfo struct {
	Address             string          `json:"address"`
	Received            decimal.Decimal `json:"received"`
	Sent                decimal.Decimal `json:"sent"`
	Balance             decimal.Decimal `json:"balance"`
	TxCount             int64           `json:"tx_count"`
	UnconfirmedReceived decimal.Decimal `json:"unconfirmed_received"`
	UnconfirmedSent     decimal.Decimal `json:"unconfirmed_sent:`
	UnconfirmedTxCount  int64           `json:"unconfirmed_tx_count"`
	UnspentTxCount      int64           `json:"unspent_tx_count"`
	FirstTx             string          `json:"first_tx"`
	LastTx              string          `json:"last_tx"`
}

type UnSpentsUTXO struct {
	TxId         string          `json:"txid"`
	Vout         int64           `json:"vout"`
	ScriptPubKey string          `json:"scriptPubKey"`
	Amount       decimal.Decimal `json:"amount"`
	Height       int64           `json:"height:`
}

//*********************************************** BTC UTXO info ***********************************************
