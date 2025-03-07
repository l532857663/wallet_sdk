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
	ETH_Mainnet     = "eth_main"
	ETH_Rinkeby     = "rinkeby"
	ETH_Sepolia     = "sepolia"
	HT_Testnet      = "ht_test"
	BSC_Testnet     = "bsc_test"
	POLYGON_Testnet = "polygon_test"
	BASE_Mainnet    = "base_main"
	BASE_Sepolia    = "base_sepolia"
	// TRON系网络
	TRX_Nile = "trx_nile"
	// BTC系网络
	BTC_Mainnet = "btc_main"
	BTC_Testnet = "btc_test"
	BTC_RegTest = "btc_reg"
	// SOLANA系网络

	// chain Symbol
	MainCoinBTC     = "BTC"
	MainCoinETH     = "ETH"
	MainCoinBSC     = "BSC"
	MainCoinPolygon = "MATIC"
	MainCoinTRON    = "TRX"
	MainCoinSOLANA  = "SOL"
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

	// 网络选择器
	ChainCombo = []string{
		ETH_Rinkeby, ETH_Sepolia, BASE_Sepolia,
		BTC_Mainnet, BTC_Testnet, BTC_RegTest,
	}

	// chain Type
	ChainRelationMap = map[string]string{
		MainCoinBTC:     MainCoinBTC,
		MainCoinETH:     MainCoinETH,
		MainCoinBSC:     MainCoinETH,
		MainCoinPolygon: MainCoinETH,
		MainCoinTRON:    MainCoinTRON,
	}
)

var (
	// ETH 节点信息配置
	ETHRinkeby = client.Node{
		ChainType: MainCoinETH,
		Ip:        "eth-node.hkva-inc.com",
		Port:      8545,
		ChainId:   "4",
	}

	ETHSepolia = client.Node{
		ChainType: MainCoinETH,
		Ip:        "eth-node.hkva-inc.com",
		Port:      8545,
		ChainId:   "11155111",
	}

	HTTestnet = client.Node{
		ChainType: MainCoinETH,
		Ip:        "https://http-testnet.hecochain.com",
		ChainId:   "256",
	}

	BSCTestnet = client.Node{
		ChainType: MainCoinETH,
		Ip:        "https://data-seed-prebsc-1-s1.binance.org:8545/",
		ChainId:   "97",
	}

	POLYGONTestnet = client.Node{
		ChainType: MainCoinETH,
		Ip:        "https://rpc-mumbai.matic.today",
		ChainId:   "80001",
	}

	TRXTestnet = client.Node{
		ChainType: MainCoinTRON,
		Ip:        "grpc.nile.trongrid.io",
		Port:      50051,
		ChainId:   "5",
	}

	BTCTestnet = client.Node{
		ChainType: MainCoinBTC,
		Ip:        "10.20.13.200",
		Port:      18443,
		User:      "btc",
		Password:  "btc2021",
		ChainId:   "",
		Net:       "testnet3",
	}

	BTCMainnet = client.Node{
		ChainType: MainCoinBTC,
		Ip:        "https://btc.getblock.io/09114f78-4075-46fa-a11f-d1678739f988/mainnet/",
		Port:      0,
		User:      "btc",
		Password:  "btc2021",
		ChainId:   "",
		Net:       "mainnet",
	}

	BTCRegtest = client.Node{
		ChainType: MainCoinBTC,
		Ip:        "10.20.13.200",
		Port:      18332,
		User:      "btc",
		Password:  "btc2021",
		ChainId:   "",
		Net:       "regtest",
	}

	BASESMainnet = client.Node{
		ChainType: MainCoinETH,
		Ip:        "https://mainnet.base.org",
		Port:      0,
		ChainId:   "8453",
	}

	BASESepolia = client.Node{
		ChainType: MainCoinETH,
		Ip:        "https://sepolia.base.org",
		Port:      0,
		ChainId:   "84532",
	}
)
