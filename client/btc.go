package client

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"wallet_sdk/txrules"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/minchenzz/brc20tool/pkg/btcapi/mempool"
)

type BtcClient struct {
	Client        *rpcclient.Client
	MempoolClient *mempool.MempoolClient
	Params        *chaincfg.Params
	Node          *Node
	PsbtUpdater   *psbt.Updater
	Placeholder   string
}

// BTC节点
func NewBtcClient(conf *Node) (*BtcClient, error) {
	var (
		url     string
		isHttps bool
	)
	if conf.Port != 0 {
		url = fmt.Sprintf("%s:%d", conf.Ip, conf.Port)
	} else {
		url = conf.Ip
	}

	// 某些https节点配置需要做一些特殊处理
	if strings.HasPrefix(url, "https://") {
		isHttps = true
		url = strings.TrimPrefix(url, "https://")
	}

	connCfg := &rpcclient.ConnConfig{
		Host:         url,
		User:         conf.User,
		Pass:         conf.Password,
		HTTPPostMode: true,     // Bitcoin core only supports HTTP POST mode
		DisableTLS:   !isHttps, // Bitcoin core does not provide TLS by default
	}
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, err
	}
	node := &BtcClient{
		Client: client,
		Node:   conf,
	}
	switch strings.ToUpper(conf.Net) {
	case BtcNodeNetMain:
		node.Params = &chaincfg.MainNetParams
		node.Placeholder = "bc1pfzl0rw44mkgevdauhrtzy5kdztjezyq0rnfqfppzxtnrwzdj553qm2vsxf"
	case BtcNodeNetTestNet3:
		node.Params = &chaincfg.TestNet3Params
		node.Placeholder = "tb1pg0uc7ujx6rplw4wj73etg505jh49k63s7wc3kyngf73ze7ffue4skru6ld"
	case BtcNodeNetRegTest:
		node.Params = &chaincfg.RegressionNetParams
	default:
		node.Params = &chaincfg.Params{}
	}
	// 初始化外部服务
	node.MempoolClient = mempool.NewClient(node.Params)
	return node, nil
}

// 查询地址余额
func (c *BtcClient) GetBalance(addr, state string) (*big.Int, error) {
	utxoList := c.getAddressUTXO(addr)
	balance := big.NewInt(0)
	for _, utxo := range utxoList {
		balance.Add(balance, utxo.RawAmount)
	}
	return balance, nil
}

// 查询地址UTXO列表
func (c *BtcClient) GetAddressUTXO(addr, state string) (interface{}, error) {
	return c.getAddressUTXO(addr), nil
}

// 根据交易HASH查交易详情
func (c *BtcClient) GetTransactionByHash(txHash string) (interface{}, bool, error) {
	h, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return nil, false, err
	}
	result, err := c.Client.GetRawTransactionVerbose(h)
	// 查询块信息
	return result, false, err
}

// 查询建议手续
func (c *BtcClient) SuggestGasPrice() *big.Int {
	// 查询链上数据
	gasPrice := big.NewInt(0)
	estimateFee, err := c.getFeePerKB(BtcEstimateSmartFeeConfirmBlock, BtcMaxFeePerKb)
	if err != nil {
		fmt.Printf("SuggestGasPrice, GetFeePerKb fatal, %+v", err)
		return gasPrice
	}
	//计算单kb手续费 单位satoshi
	satoshi, err := btcutil.NewAmount(estimateFee)
	if nil != err {
		fmt.Printf("SuggestGasPrice, NewAmount fatal, %+v", err)
		return gasPrice
	}
	gasPrice = big.NewInt(int64(satoshi))
	return gasPrice
}

func (c *BtcClient) getFeePerKB(nblocks int64, btcMaxFeePerKb float64) (float64, error) {
	feeInfo, err := c.Client.EstimateSmartFee(nblocks, &btcjson.EstimateModeConservative)
	if err != nil {
		return 0, fmt.Errorf("getFeePerKB fatal, " + err.Error())
	}
	feePerKb := *feeInfo.FeeRate
	if feePerKb > btcMaxFeePerKb {
		feePerKb = btcMaxFeePerKb
	}
	return feePerKb, nil
}

// 查询最新区块高度
func (c *BtcClient) GetBlockHeight() (int64, error) {
	return c.Client.GetBlockCount()
}

// 构建交易
func (c *BtcClient) BuildTransferInfo(fromAddr, toAddr, contract, amount, gasPrice, nonce string) (interface{}, error) {
	// to地址数据处理
	detail := GetToAddrDetail(toAddr, amount)
	toAddrList := []*ToAddrDetail{
		detail,
	}
	// 手续费
	gas := BtcToSatoshi(gasPrice)
	// 查询FROM地址的UTXO
	unspendUTXOList := c.getAddressUTXO(fromAddr)
	if len(unspendUTXOList) < 1 {
		return nil, fmt.Errorf("BuildTransaction fatal, from address not unspend utxo")
	}
	// 最终使用UTXO
	var useUTXOList []*UnspendUTXOList
	// 排序未花费得UTXO
	DescSortUnspendUTXO(unspendUTXOList)
	inAmount := big.NewInt(0)
	outAmount := big.NewInt(0).Add(detail.RawAmount, gas)
	for _, info := range unspendUTXOList {
		if info.Amount < 700 {
			continue
		}
		useUTXOList = append(useUTXOList, info)
		inAmount.Add(inAmount, info.RawAmount)
		if inAmount.Cmp(outAmount) > 0 {
			break
		}
	}

	apiTx, err := c.genBtcTransaction(useUTXOList, toAddrList, gas, toAddr)
	if nil != err {
		return nil, fmt.Errorf("BuildTransaction fatal, " + err.Error())
	}
	txObj := &BtcTransferInfo{
		ApiTx:    apiTx,
		UTXOList: useUTXOList,
	}
	return txObj, nil
}

func (c *BtcClient) getAddressUTXO(address string) []*UnspendUTXOList {
	var res []*UnspendUTXOList
	// 使用外部服务
	addr, err := btcutil.DecodeAddress(address, c.Params)
	if err != nil {
		fmt.Printf("invalid recipet address: %v", err)
		return nil
	}
	// 查询未花费的UTXO列表
	unspendList, err := c.MempoolClient.ListUnspent(addr)
	if err != nil {
		fmt.Printf("GetListUnspent error: %+v", err)
		return nil
	}
	if len(unspendList) == 0 {
		fmt.Printf("no utxo for %v", addr)
		return nil
	}
	// ScriptPubKey
	spk, err := txscript.PayToAddrScript(addr)
	if err != nil {
		fmt.Printf("PayToAddrScript err: %v", err)
		return nil
	}
	// 格式化
	for _, unspend := range unspendList {
		amount := unspend.Output.Value
		tmp := &UnspendUTXOList{
			TxHash:       unspend.Outpoint.Hash.String(),
			ScriptPubKey: hex.EncodeToString(spk),
			Vout:         unspend.Outpoint.Index,
			Amount:       amount,
			RawAmount:    big.NewInt(amount),
		}
		res = append(res, tmp)
	}
	return res
}

func (c *BtcClient) genBtcTransaction(unSpendUTXOList []*UnspendUTXOList, toAddrList []*ToAddrDetail, gasPrice *big.Int, changeAddr string) (*wire.MsgTx, error) {
	retApiTx := wire.NewMsgTx(wire.TxVersion)
	// 计算 UTXO
	var inAmount btcutil.Amount
	for _, rti := range unSpendUTXOList {
		txHash, err := chainhash.NewHashFromStr(rti.TxHash)
		if err != nil {
			return nil, fmt.Errorf("genBtcTransaction NewHashFromStr fatal, %v, tx hash: %s", err, rti.TxHash)
		}
		satoshi := btcutil.Amount(rti.RawAmount.Int64())
		inAmount += satoshi
		prevOut := wire.NewOutPoint(txHash, rti.Vout)
		retApiTx.AddTxIn(wire.NewTxIn(prevOut, nil, nil))
	}
	fmt.Printf("inAmount: %+v\n", inAmount)

	// 接收地址
	var outAmount btcutil.Amount
	for _, detail := range toAddrList {
		toAddr := detail.Address
		satoshi := btcutil.Amount(detail.RawAmount.Int64())
		if satoshi <= 0 {
			continue
		}
		fmt.Printf("toAddr: %+v, value: %+v\n", toAddr, satoshi)
		outAmount += satoshi
		// Decode the recipent address.
		pkScript, err := NewPubkeyHash(toAddr, c.Params)
		if err != nil {
			return nil, fmt.Errorf("genBtcTransaction NewPubkeyHash fatal, " + err.Error())
		}
		txOut := wire.NewTxOut(int64(satoshi), pkScript)
		err = txrules.CheckOutput(txOut, txrules.DefaultRelayFeePerKb)
		if err != nil {
			return nil, fmt.Errorf("genBtcTransaction CheckOutput fatal, " + err.Error())
		}
		retApiTx.AddTxOut(txOut)
	}

	if outAmount >= inAmount {
		return nil, fmt.Errorf("genBtcTransaction fatal, outAmount: %d >= inAmount: %d", int64(outAmount), int(inAmount))
	}
	// gasPrice
	relayFeePerKb := btcutil.Amount(gasPrice.Int64())
	if relayFeePerKb == 0 {
		return nil, fmt.Errorf("genBtcTransaction fatal, feePerKB ni zero")
	}

	//进行找零 ScriptPubKey
	pkScript, err := NewPubkeyHash(changeAddr, c.Params)
	if err != nil {
		return nil, fmt.Errorf("genBtcTransaction txscript.PayToAddrScript fatal, %s, addr: %s", err.Error(), changeAddr)
	}
	txOut := wire.NewTxOut(0, pkScript)
	retApiTx.AddTxOut(txOut)

	// 交易内容大小
	preTxSize := PreCalculateSerializeSize(retApiTx)
	fmt.Printf("preTxSize: %+v\n", preTxSize)
	// 预估手续费
	finalFee := txrules.FeeForSerializeSize(relayFeePerKb, preTxSize)
	fmt.Printf("finalFee: %+v\n", finalFee)

	//btc 每条交易限制不超过100Kb
	if preTxSize >= BtcMaxTransactionByteSizeKB*1000 {
		return nil, fmt.Errorf("genBtcTransaction fatal, tx pre-calculate size is more than 100Kb")
	}
	//找零金额去掉手续费
	changeAmount := inAmount - outAmount - finalFee
	fmt.Printf("changeAmount: %+v\n", changeAmount)
	//一个utxo使用的时候会产生148个字节的手续费，如果找零金额 低于下面值就不找零，直接当手续费了
	minChange := BtcMinChangeByte * relayFeePerKb / 1000
	//此时不够说明找零没有意义了，直接将找零当成手续费就可以了
	if changeAmount < minChange && changeAmount >= 0 {
		retApiTx.TxOut = retApiTx.TxOut[:len(retApiTx.TxOut)-1]
	} else if changeAmount < 0 {
		return nil, fmt.Errorf("genBtcTransaction fatal, gasFee is not enough")
	} else {
		retApiTx.TxOut[len(retApiTx.TxOut)-1].Value = int64(changeAmount)

		//检查找零是否dust，dust就直接不找零了
		isDust := txrules.IsDustOutput(retApiTx.TxOut[len(retApiTx.TxOut)-1], txrules.DefaultRelayFeePerKb)
		if isDust {
			retApiTx.TxOut = retApiTx.TxOut[:len(retApiTx.TxOut)-1]
			fmt.Printf("genBtcTransaction change amount dust! retApiTx: %+v\n", retApiTx)
		}
	}
	return retApiTx, nil
}

// SignTransferToRaw(txObj, hexPrivateKey string) (string, error)
func (c *BtcClient) SignTransferToRaw(txObj, hexPrivateKey string) (string, error) {
	txInfo := &BtcTransferInfo{}
	err := json.Unmarshal([]byte(txObj), txInfo)
	if err != nil {
		return "", err
	}
	apiTx := txInfo.ApiTx
	fmt.Printf("wch----- apiTx: %+v\n", apiTx)
	// 解析私钥
	wif, err := btcutil.DecodeWIF(hexPrivateKey)
	if err != nil {
		return "", fmt.Errorf("SignTx DecodeWIF fatal, " + err.Error())
	}
	if !wif.IsForNet(c.Params) {
		return "", fmt.Errorf("SignTx IsForNet not matched")
	}
	prevOutFetcher := txscript.NewMultiPrevOutFetcher(nil)
	var privateKeys []*btcec.PrivateKey
	for i := 0; i < len(apiTx.TxIn); i++ {
		input := apiTx.TxIn[i]
		utxoInfo := txInfo.UTXOList[i]
		outPoint := &input.PreviousOutPoint
		pkScript, err := hex.DecodeString(utxoInfo.ScriptPubKey)
		if err != nil {
			return "", err
		}
		txOut := wire.NewTxOut(utxoInfo.Amount, pkScript)
		prevOutFetcher.AddPrevOut(*outPoint, txOut)
		privateKeys = append(privateKeys, wif.PrivKey)
	}

	if err := Sign(apiTx, privateKeys, prevOutFetcher); err != nil {
		return "", err
	}

	raw, err := getTxHex(apiTx)
	if err != nil {
		return "", err
	}
	fmt.Printf("apiTx txHash: %s, info: %s\n", apiTx.TxHash(), raw)
	return raw, nil
}

func Sign(tx *wire.MsgTx, privateKeys []*btcec.PrivateKey, prevOutFetcher *txscript.MultiPrevOutFetcher) error {
	for i, in := range tx.TxIn {
		prevOut := prevOutFetcher.FetchPrevOutput(in.PreviousOutPoint)
		txSigHashes := txscript.NewTxSigHashes(tx, prevOutFetcher)
		privKey := privateKeys[i]
		if txscript.IsPayToTaproot(prevOut.PkScript) {
			witness, err := txscript.TaprootWitnessSignature(tx, txSigHashes, i, prevOut.Value, prevOut.PkScript, txscript.SigHashDefault, privKey)
			if err != nil {
				return err
			}
			in.Witness = witness
		} else if txscript.IsPayToPubKeyHash(prevOut.PkScript) {
			sigScript, err := txscript.SignatureScript(tx, i, prevOut.PkScript, txscript.SigHashAll, privKey, true)
			if err != nil {
				return err
			}
			in.SignatureScript = sigScript
		} else {
			pubKeyBytes := privKey.PubKey().SerializeCompressed()
			script, err := PayToPubKeyHashScript(btcutil.Hash160(pubKeyBytes))
			if err != nil {
				return err
			}
			amount := prevOut.Value
			witness, err := txscript.WitnessSignature(tx, txSigHashes, i, amount, script, txscript.SigHashAll, privKey, true)
			if err != nil {
				return err
			}
			in.Witness = witness

			if txscript.IsPayToScriptHash(prevOut.PkScript) {
				redeemScript, err := PayToWitnessPubKeyHashScript(btcutil.Hash160(pubKeyBytes))
				if err != nil {
					return err
				}
				in.SignatureScript = append([]byte{byte(len(redeemScript))}, redeemScript...)
			}
		}
	}

	return nil
}

func (c *BtcClient) sign(apiTx *wire.MsgTx, privateKey string, utxoIdx int, utxoSciptPubKey []byte) (string, error) {
	txIn := apiTx.TxIn[utxoIdx]
	if nil == txIn {
		fmt.Printf("btc sign TxIn err! TxIn: %+v utxoIdx: %+v\n", apiTx.TxIn, utxoIdx)
		return "", fmt.Errorf("btc SignTx txIn is nil")
	}
	// 解析私钥
	wif, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "", fmt.Errorf("SignTx DecodeWIF fatal, " + err.Error())
	}
	if !wif.IsForNet(c.Params) {
		return "", fmt.Errorf("SignTx IsForNet not matched")
	}
	getKey := txscript.KeyClosure(func(addr btcutil.Address) (*btcec.PrivateKey, bool, error) {
		return wif.PrivKey, wif.CompressPubKey, nil
	})

	getScript := txscript.ScriptClosure(func(addr btcutil.Address) ([]byte, error) {
		pubKey, err := btcutil.NewAddressPubKey(wif.SerializePubKey(), c.Params)
		if err != nil {
			fmt.Printf("btc sign NewAddressPubKey err! err: %+v\n", err)
			return nil, err
		}
		return txscript.MultiSigScript([]*btcutil.AddressPubKey{pubKey}, 1)
	})

	txIn.SignatureScript, err = txscript.SignTxOutput(c.Params, apiTx, utxoIdx, utxoSciptPubKey,
		txscript.SigHashAll, getKey, getScript, nil)
	if err != nil {
		fmt.Printf("btc sign SignTxOutput err! apiTx: %+v utxoIdx: %+v err: %+v\n", apiTx, utxoIdx, err)
		return "", fmt.Errorf("SignTx SignTxOutput fatal, " + err.Error())
	}

	vm, err := txscript.NewEngine(utxoSciptPubKey, apiTx, utxoIdx, txscript.StandardVerifyFlags, nil, nil, 0, nil)
	if err != nil {
		fmt.Printf("btc sign NewEngine err! apiTx: %+v utxoIdx: %+v err: %+v\n", apiTx, utxoIdx, err)
		return "", fmt.Errorf("SignTx NewEngine fatal, " + err.Error())
	}
	if err = vm.Execute(); err != nil {
		return "", fmt.Errorf("SignTx Execute fatal, " + err.Error())
	}
	fmt.Println("btc sign end")
	return "", nil
}

// SendRawTransaction sends tx to node.
// SendRawTransaction(hexTx string) (string, error)
func (c *BtcClient) SendRawTransaction(hexTx string) (string, error) {
	method := "sendrawtransaction"
	var txHash string
	marshalledParam, err := json.Marshal(hexTx)
	if err != nil {
		err = fmt.Errorf("SendRawTransaction json.Marshal error: %+v", err)
		return txHash, err
	}
	rawMessage := json.RawMessage(marshalledParam)
	result, err := c.Client.RawRequest(method, []json.RawMessage{rawMessage})
	if err != nil {
		err = fmt.Errorf("SendRawTransaction RawRequest error: %+v", err)
		return txHash, err
	}
	err = json.Unmarshal(result, &txHash)
	if err != nil {
		err = fmt.Errorf("SendRawTransaction json.Unmarshal error: %+v", err)
		return txHash, err
	}
	return txHash, nil
}

// 构建交易
func (c *BtcClient) BuildTransferInfoByList(unSpendUTXOList []*UnspendUTXOList, toAddrList []*ToAddrDetail, gasPrice, changeAddr string) (interface{}, error) {
	// 手续费
	gas := BtcToSatoshi(gasPrice)

	apiTx, err := c.genBtcTransaction(unSpendUTXOList, toAddrList, gas, changeAddr)
	if nil != err {
		return nil, fmt.Errorf("BuildTransaction fatal, " + err.Error())
	}
	txObj := &BtcTransferInfo{
		ApiTx:    apiTx,
		UTXOList: unSpendUTXOList,
	}
	return txObj, nil
}

// 多个地址的签名出账
func (c *BtcClient) SignListAndSendTransfer(txObj string, hexPrivateKeys []string) (string, error) {
	txInfo := &BtcTransferInfo{}
	err := json.Unmarshal([]byte(txObj), txInfo)
	if err != nil {
		return "", err
	}
	apiTx := txInfo.ApiTx
	for idx, rti := range txInfo.UTXOList {
		prevOutScript, err := hex.DecodeString(rti.ScriptPubKey)
		if err != nil {
			fmt.Printf("invalid script key error: %+v\n", err)
			return "", err
		}
		_, err = c.sign(apiTx, hexPrivateKeys[idx], idx, prevOutScript)
		if err != nil {
			fmt.Printf("Sign err: %+v\n", err)
			return "", err
		}
	}
	// 签名
	var buf bytes.Buffer
	buf.Grow(hex.EncodedLen(apiTx.SerializeSize()))
	if err := apiTx.Serialize(hex.NewEncoder(&buf)); err != nil {
		return "", err
	}
	fmt.Printf("apiTx info: %+v\n", buf.String())

	txHash, err := c.Client.SendRawTransaction(apiTx, false)
	if nil != err {
		return "", fmt.Errorf("Broadcast SendRawTransaction fatal, " + err.Error())
	}
	return txHash.String(), nil
}

// 查询使用的节点信息
func (c *BtcClient) GetNodeInfo() (*Node, error) {
	return c.Node, nil
}

// 查询地址是否是Taproot类型
func (c *BtcClient) GetAddressIsTaproot(address string) bool {
	addr, err := btcutil.DecodeAddress(address, c.Params)
	if err != nil {
		fmt.Printf("btcutil.DecodeAddress Error: %+v\n", err)
		return false
	}
	switch t := addr.(type) {
	case *btcutil.AddressPubKeyHash:
		fmt.Println("wch---- PKH")
	case *btcutil.AddressScriptHash:
		fmt.Println("wch---- SH")
	case *btcutil.AddressWitnessPubKeyHash:
		fmt.Println("wch---- WPKH")
	case *btcutil.AddressWitnessScriptHash:
		fmt.Println("wch---- WSH")
	case *btcutil.AddressTaproot:
		fmt.Println("wch---- T")
		return true
	default:
		fmt.Printf("invalid recipet address type: %v\n", t)
	}
	return false
}

// 关闭链接
func (node *BtcClient) Close() {
	if node.Client != nil {
		node.Client.Shutdown()
	}
}

/** 链暂不支持的方法
 */
// 查询地址代币余额
func (c *BtcClient) GetBalanceByContract(addr, contractAddr string) (*big.Int, error) {
	return nil, fmt.Errorf("This method is not supported yet!")
}

// 查询合约精度
func (c *BtcClient) GetDecimals(contractAddr string) (*big.Int, error) {
	return nil, fmt.Errorf("This method is not supported yet!")
}

// 查询合约symbol
func (c *BtcClient) GetSymbol(contractAddr string) (string, error) {
	return "", fmt.Errorf("This method is not supported yet!")
}

// 查询地址的nonce
func (c *BtcClient) GetNonce(addr, param string) (uint64, error) {
	return 0, fmt.Errorf("This method is not supported yet!")
}

// 查询chain_id
func (c *BtcClient) ChainID() (*big.Int, error) {
	return nil, fmt.Errorf("This method is not supported yet!")
}

// 构建合约调用
func (c *BtcClient) BuildContractInfo(contract, abiContent, gasPrice, nonce, params string, args ...interface{}) (interface{}, error) {
	return nil, fmt.Errorf("This method is not supported yet!")
}
