package client

import (
	"github.com/shopspring/decimal"
	"math/big"

	"github.com/btcsuite/btcd/txscript"
)

// 节点配置信息结构体
type Node struct {
	ChainType string // 链类型
	Ip        string // IP地址 或 域名
	Port      uint64 // 端口号
	User      string // 用户名
	Password  string // 密码
	ChainId   string // 链配置
	Net       string // 网络类型(BTC等链使用)
}

var FreeNodeMap map[string]*Node // 自定义节点信息

// BTC系列的UnspendUTXOList
type UnspendUTXOList struct {
	TxHash       string          `json:"txid"`
	ScriptPubKey string          `json:"scriptPubKey"`
	Vout         uint32          `json:"vout"`
	Amount       decimal.Decimal `json:"amount"`
	RawAmount    *big.Int
}

// to地址信息
type ToAddrDetail struct {
	Address   string
	Amount    string
	RawAmount *big.Int
}

type Input struct {
	TxId              string
	VOut              uint32
	Sequence          uint32
	Amount            int64
	Address           string
	PrivateKey        string
	NonWitnessUtxo    string
	MasterFingerprint uint32
	DerivationPath    string
	PublicKey         string
}

type Output struct {
	Address           string
	Amount            int64
	IsChange          bool
	MasterFingerprint uint32
	DerivationPath    string
	PublicKey         string
}

type UtxoType int

const (
	NonWitness UtxoType = 1
	Witness    UtxoType = 2
)

type InputUtxo struct {
	UtxoType            UtxoType             `json:"utxo_type"`
	SighashType         txscript.SigHashType `json:"sighash_type"`
	NonWitnessUtxo      string               `json:"non_witness_utxo"`       //
	WitnessUtxoPkScript string               `json:"witness_utxo_pk_script"` //
	WitnessUtxoAmount   uint64               `json:"witness_utxo_amount"`    //
	Index               int                  `json:"index"`
}

type InputSigner struct {
	UtxoType    UtxoType             `json:"utxo_type"`
	SighashType txscript.SigHashType `json:"sighash_type"`
	//Sig   string `json:"sig"`
	Pri   string `json:"pri"`
	Pub   string `json:"pub"`
	Index int    `json:"index"`
}
