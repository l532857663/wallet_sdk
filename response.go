package wallet_sdk

// 返回结构
type Response struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

// 使用通用返回结果
type CommonResp struct {
	Status Response `json:"status"`
	Data   string   `json:"data"`
}

// 使用助记词生成账户信息
type AccountInfoResp struct {
	Status Response     `json:"status"`
	Data   *AccountInfo `json:"data"`
}

// 查询合约信息
type ContractInfoResp struct {
	Status Response      `json:"status"`
	Data   *ContractInfo `json:"data"`
}

// 查询交易信息
type TransactionInfoResp struct {
	Status Response         `json:"status"`
	Data   *TransactionInfo `json:"data"`
}

// 建议手续费信息
type TransactionGasPriceResp struct {
	Status Response      `json:"status"`
	Data   *GasPriceInfo `json:"data"`
}

// 查询地址的UTXOList
type AddressUTXOListResp struct {
	Status Response    `json:"status"`
	Data   interface{} `json:"data"`
}
