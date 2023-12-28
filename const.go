package wallet_sdk

import (
	"wallet_sdk/client"

	"github.com/shopspring/decimal"
)

// 常用参数
const (
	// 返回结果参数
	RES_CODE_SUCCESS = 0
	RES_CODE_FAILED  = 1

	// jsonRpc state
	StateEarliest = "earliest"
	StateLatest   = "latest"
	StatePending  = "pending"
	// contract abi params
	ContractDecimals = "decimals"
	ContractSymbol   = "symbol"

	// ETH系网络
	ETH_Rinkeby     = "rinkeby"
	ETH_Sepolia     = "sepolia"
	HT_Testnet      = "ht_test"
	BSC_Testnet     = "bsc_test"
	POLYGON_Testnet = "polygon_test"
	TRX_Nile        = "trx_nile"
	BTC_Testnet     = "btc_test"
	BTC_Mainnet     = "btc_main"

	// chain Symbol
	MainCoinEth = "ETH"
	MainCoinBSC = "BSC"

	// chain Type
	ChainRelationForBTC = "BTC"
	ChainRelationForETH = "ETH"
	ChainRelationForTRX = "TRX"
)

var (
	// 返回结果状态
	ResSuccess = Response{
		Code: RES_CODE_SUCCESS,
	}
	ResFailed = Response{
		Code: RES_CODE_FAILED,
	}

	// Gas coefficient
	GasFast, _    = decimal.NewFromString("1.5")
	GasHigh, _    = decimal.NewFromString("1.3")
	GasAverage, _ = decimal.NewFromString("1.1")
)

var (
	// ETH 节点信息配置
	ETHRinkeby = client.Node{
		ChainType: ChainRelationForETH,
		Ip:        "192.168.10.173",
		Port:      8545,
		ChainId:   "4",
	}

	ETHSepolia = client.Node{
		ChainType: ChainRelationForETH,
		Ip:        "192.168.10.173",
		Port:      8545,
		ChainId:   "11155111",
	}

	HTTestnet = client.Node{
		ChainType: ChainRelationForETH,
		Ip:        "https://http-testnet.hecochain.com",
		ChainId:   "256",
	}

	BSCTestnet = client.Node{
		ChainType: ChainRelationForETH,
		Ip:        "https://data-seed-prebsc-1-s1.binance.org:8545/",
		ChainId:   "97",
	}

	POLYGONTestnet = client.Node{
		ChainType: ChainRelationForETH,
		Ip:        "https://rpc-mumbai.matic.today",
		ChainId:   "80001",
	}

	TRXTestnet = client.Node{
		ChainType: ChainRelationForTRX,
		Ip:        "grpc.nile.trongrid.io",
		Port:      50051,
		ChainId:   "5",
	}

	BTCTestnet = client.Node{
		ChainType: ChainRelationForBTC,
		Ip:        "128.140.73.158",
		Port:      18443,
		User:      "btc",
		Password:  "btc2021",
		ChainId:   "",
		Net:       "testnet3",
	}

	BTCMainnet = client.Node{
		ChainType: ChainRelationForBTC,
		Ip:        "https://btc.getblock.io/09114f78-4075-46fa-a11f-d1678739f988/mainnet/",
		Port:      0,
		User:      "btc",
		Password:  "btc2021",
		ChainId:   "",
		Net:       "mainnet",
	}
)
