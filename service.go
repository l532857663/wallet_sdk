package wallet_sdk

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"wallet_sdk/client"
)

type NodeService interface {
	// 查询账户信息
	GetBalance(addr, state string) (*big.Int, error)
	GetBalanceByContract(addr, contractAddr string) (*big.Int, error)
	GetTransactionByHash(txHash string) (interface{}, bool, error)
	GetAddressUTXO(addr, state string) (interface{}, error)
	GetAddressIsTaproot(addr string) bool

	// 查询合约信息
	GetDecimals(contractAddr string) (*big.Int, error)
	GetSymbol(contractAddr string) (string, error)
	GetContractInfoByFunc(contractAddr, funcName string, args ...interface{}) (interface{}, error)

	// 查询账户交易相关信息
	SuggestGasPrice() *big.Int
	GetNonce(addr, param string) (uint64, error)

	// 发送交易相关
	BuildTransferInfo(fromAddr, toAddr, contract, amount, gasPrice, nonce string) (interface{}, error)
	BuildContractInfo(contract, abiContent, gasPrice, nonce, params string, args ...interface{}) (interface{}, error)
	SignTransferToRaw(txObj, hexPrivateKey string) (string, error)
	SendRawTransaction(hexTx string) (string, error)
	BuildPSBTransfer(ins []client.Input, outs []client.Output) (interface{}, error)
	BuildTransferInfoByList(fromAddr []*client.UnspendUTXOList, toAddr []*client.ToAddrDetail, gasPrice, changeAddr string) (interface{}, error)
	SignListAndSendTransfer(txObj string, hexPrivateKeys []string) (string, error)
	GenerateSignedListingPSBTBase64(ins *client.Input, outs *client.Output) (interface{}, error)

	ChainID() (*big.Int, error)
	GetBlockHeight() (int64, error)

	// 查询节点信息
	GetNodeInfo() (*client.Node, error)
	Close()
}

func NewNodeService(chainName string) (NodeService, error) {
	// 链接节点配置信息
	switch chainName {
	case ETH_Rinkeby:
		cli, err := client.NewEthClient(&ETHRinkeby)
		if err != nil {
			return nil, err
		}
		return cli, nil
	case ETH_Sepolia:
		cli, err := client.NewEthClient(&ETHSepolia)
		if err != nil {
			return nil, err
		}
		return cli, nil
	case BSC_Testnet:
		cli, err := client.NewEthClient(&BSCTestnet)
		if err != nil {
			return nil, err
		}
		return cli, nil
	case POLYGON_Testnet:
		cli, err := client.NewEthClient(&POLYGONTestnet)
		if err != nil {
			return nil, err
		}
		return cli, nil
	case TRX_Nile:
		cli, err := client.NewTronClient(&TRXTestnet)
		if err != nil {
			return nil, err
		}
		return cli, nil
	case BTC_Testnet:
		cli, err := client.NewBtcClient(&BTCTestnet)
		if err != nil {
			return nil, err
		}
		return cli, nil
	case BTC_Mainnet:
		cli, err := client.NewBtcClient(&BTCMainnet)
		if err != nil {
			return nil, err
		}
		return cli, nil
	case BTC_RegTest:
		cli, err := client.NewBtcClient(&BTCRegtest)
		if err != nil {
			return nil, err
		}
		return cli, nil
	case BASE_Mainnet:
		cli, err := client.NewEthClient(&BASESMainnet)
		if err != nil {
			return nil, err
		}
		return cli, nil
	case BASE_Sepolia:
		cli, err := client.NewEthClient(&BASESepolia)
		if err != nil {
			return nil, err
		}
		return cli, nil
	}
	cli, err := GetFreeNodeSerivce(chainName)
	if err != nil {
		return nil, err
	}
	return cli, nil

}

func SetNodeInfo(name, chainType, ip, portStr, user, password, chainId, net string) {
	if client.FreeNodeMap == nil {
		client.FreeNodeMap = make(map[string]*client.Node)
	}
	port, _ := strconv.ParseUint(portStr, 0, 64)
	node := &client.Node{
		ChainType: chainType,
		Ip:        ip,
		Port:      port,
		User:      user,
		Password:  password,
		ChainId:   chainId,
		Net:       net,
	}
	client.FreeNodeMap[name] = node
}

func GetFreeNodeSerivce(name string) (NodeService, error) {
	node, ok := client.FreeNodeMap[name]
	if !ok {
		return nil, fmt.Errorf("Not get node[%s] info!", name)
	}
	chainType := strings.ToUpper(node.ChainType)
	switch chainType {
	case "ETH", "ETHEREUM":
		cli, err := client.NewEthClient(node)
		if err != nil {
			return nil, err
		}
		return cli, nil
	case "TRX", "TRON":
		cli, err := client.NewTronClient(node)
		if err != nil {
			return nil, err
		}
		return cli, nil
	case "BTC", "BITCOIN":
		cli, err := client.NewBtcClient(node)
		if err != nil {
			return nil, err
		}
		return cli, nil
	}
	return nil, fmt.Errorf("The chain[%s] is not support!", node.ChainType)
}
