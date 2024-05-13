package wallet_sdk

// 账户信息
type AccountInfo struct {
	Address     string       `json:"address"`
	PrivateKey  string       `json:"privateKay"`
	PublicKey   string       `json:"publicKey"`
	BtcAddrList []BtcAddress `json:"btcList"`
}

// 合约信息
type ContractInfo struct {
	Symbol   string `json:"symbol"`   // 代币名称
	Decimals string `json:"decimals"` // 精度
}

// 交易信息
type TransactionInfo struct {
	TxInfo    interface{} `json:"txInfo"` // 交易信息
	IsPending bool        `json:"isPending"`
}

// GasPrice建议
type GasPriceInfo struct {
	Fast    string `json:"fast"`    // 极快
	High    string `json:"high"`    // 快
	Average string `json:"average"` // 建议手续费
	Low     string `json:"low"`     // 低
}

// 选择使用的UTXO
type ChooseUTXO struct {
	TxHash string
	Vout   uint32
}

// BTC地址及类型
type BtcAddress struct {
	Address     string `json:"address"`
	PrivateKey  string `json:"privateKay"`
	AddressType string `json:"addressType"`
}
