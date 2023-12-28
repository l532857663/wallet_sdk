package client

import (
	"time"

	"github.com/btcsuite/btcd/txscript"
	"github.com/shopspring/decimal"
)

const (
	MainCoinEth = "ETH"

	EthGasLimit   uint64 = 21000
	Erc20GasLimit uint64 = 80000
	Trc20FeeLimit int64  = 10000000

	EthBase       int64 = 1000000000
	GweiLength    int32 = 9
	BtcBase       int64 = 100000000
	SatoshiLength int32 = 8
	TrxBase       int64 = 1000000
)

var (
	EthBaseDecimal  = decimal.NewFromInt(EthBase)
	TransferEventID = []byte{0xdd, 0xf2, 0x52, 0xad, 0x1b, 0xe2, 0xc8, 0x9b, 0x69, 0xc2, 0xb0, 0x68, 0xfc, 0x37, 0x8d, 0xaa, 0x95, 0x2b, 0xa7, 0xf1, 0x63, 0xc4, 0xa1, 0x16, 0x28, 0xf5, 0x5a, 0x4d, 0xf5, 0x23, 0xb3, 0xef}

	BtcBaseDecimal = decimal.NewFromInt(BtcBase)
)

const (
	// BtcMaxFeePerKb btc每kb最大fee
	BtcMaxFeePerKb = float64(0.0001)

	// BtcMaxSigScriptByteSize spendSize is the largest number of bytes of a sigScript
	// which spends a p2pkh output: OP_DATA_73 <sig> OP_DATA_33 <pubkey>
	// https://vimsky.com/zh-tw/examples/detail/golang-ex-github.com.btcsuite.btcd.wire-MsgTx-SerializeSize-method.html
	BtcMaxSigScriptByteSize = 1 + 73 + 1 + 33

	// BtcMergeMaxWaitDuration btc合并订单最大等待时间，默认10分钟
	// 配置在builder.apollo
	BtcMergeMaxWaitDuration = 10 * time.Minute

	// BtcEstimateSmartFeeConfirmBlock btc取费率区块数，默认值
	// 解释：预估这笔交易在经过几个区块后被打包
	BtcEstimateSmartFeeConfirmBlock = int64(2)

	// BtcMaxTransactionByteSizeKB btc每条交易最大字节数（单位KB, 1000字节）
	BtcMaxTransactionByteSizeKB = 100

	// BtcMinChangeByte btc最小找零金额对应的字节数
	// 即如果产生找零，这个找零被花掉的成本都大于这个找零，就不划算了
	// 写100是因为把这个判断标准降低点，一个txin占用148字节
	BtcMinChangeByte = 100
)

var (
	// btc 网络类型, 全部大写
	BtcNodeNetMain     = "MAINNET"
	BtcNodeNetTestNet3 = "TESTNET3"
	BtcNodeNetRegTest  = "REGTEST"
)

var (
	InUtxos = []InputUtxo{
		{
			UtxoType:            Witness,
			SighashType:         txscript.SigHashSingle,
			NonWitnessUtxo:      "",
			WitnessUtxoPkScript: "0014d8c9cf87df6269a9962023a57f18b93d1e4417fa",
			WitnessUtxoAmount:   1000,
			Index:               0,
		},
		{
			UtxoType:            NonWitness,
			SighashType:         txscript.SigHashSingle,
			NonWitnessUtxo:      "0200000002794d4682943dd8ef92b3f4c84a1d4b377d70f0e179dc8e87884af88bf56294d3010000006a473044022052abb42fa2ba9ea90b07cdbac64ca3a52e9b4bd364f0145ee3df8e54c1127ef102203b1cd23bf50079d7491af49e4f64a489afeadebebc8a0f85313e4ef77b636cf803210339184fdced859e6743d793ed45c52e3fba739fda2e86137b70e1dd1f13974d1dffffffffa943093bbdd97655eeb0531056aff88148279b4ea1043131066196e71f868887010000006b4830450221008f138758d9690887871d14c96488cd0e8dd982c4d30bc745109996db5b749af902204dd14ebd7a6bfaaa3af0a463bfc08dcf262205cc60ad2685385056f3acc864a203210339184fdced859e6743d793ed45c52e3fba739fda2e86137b70e1dd1f13974d1dffffffff0258020000000000001976a914d8c9cf87df6269a9962023a57f18b93d1e4417fa88acb0360000000000001976a914d8c9cf87df6269a9962023a57f18b93d1e4417fa88ac00000000",
			WitnessUtxoPkScript: "",
			WitnessUtxoAmount:   14000,
			Index:               1,
		},
	}

	inSigners1 = []InputSigner{
		{
			UtxoType:    Witness,
			SighashType: txscript.SigHashSingle,
			Pri:         "",
			Pub:         "512043f98f7246d0c3f755d2f472b451f495ea5b6a30f3b11b12684fa22cf929e66b",
			Index:       0,
		},
	}
	inSigners2 = []InputSigner{
		{
			UtxoType:    NonWitness,
			SighashType: txscript.SigHashSingle,
			Pri:         "",
			Pub:         "0339184fdced859e6743d793ed45c52e3fba739fda2e86137b70e1dd1f13974d1d",
			Index:       1,
		},
	}
)
