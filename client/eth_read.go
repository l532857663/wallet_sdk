package client

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"time"
	"wallet_sdk/utils"
)

// 查询地址余额 GetBalance(addr, state string) (*big.Int, error)
// 查询地址代币余额 GetBalanceByContract(addr, contractAddr string) (*big.Int, error)
// 查询交易信息 GetTransactionByHash(txHash string) (interface{}, bool, error)
// 查询合约精度 GetDecimals(contractAddr string) (*big.Int, error)
// 查询合约symbol GetSymbol(contractAddr string) (string, error)
// 查询地址的nonce GetNonce(addr, param string) (uint64, error)
// 查询chain_id ChainID() (*big.Int, error)
// 查询最新区块高度 GetBlockHeight() (int64, error)
// 查询建议手续费 SuggestGasPrice() *big.Int
// 查询合约信息 GetContractInfoByFunc(contractAddr, params string, args ...interface{}) (interface{}, error)

func (c *EthClient) GetBalance(addr, state string) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	account := EthAddressChange(addr)
	defer cancel()
	var result hexutil.Big
	err := c.RpcClient.CallContext(ctx, &result, "eth_getBalance", account, state)
	return (*big.Int)(&result), err
}

func (c *EthClient) GetBalanceByContract(addr, contractAddr string) (*big.Int, error) {
	result, err := c.ContractCall(contractAddr, "balanceOf", addr, contractAddr, EthAddressChange(addr))
	if err != nil {
		return nil, err
	}
	return utils.ByteTobigInt(result), nil
}

func (c *EthClient) GetTransactionByHash(txHash string) (interface{}, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	return c.TransactionByHash(ctx, HexToHash(txHash))
}

func (c *EthClient) GetDecimals(contractAddr string) (*big.Int, error) {
	result, err := c.ContractCall(contractAddr, "decimals", contractAddr, contractAddr)
	if err != nil {
		return nil, err
	}
	return utils.ByteTobigInt(result), nil
}

func (c *EthClient) GetSymbol(contractAddr string) (string, error) {
	result, err := c.ContractCall(contractAddr, "symbol", contractAddr, contractAddr)
	if err != nil {
		return "", err
	}
	// NOTE：该方法返回数据总长度为96
	if len(result) != 96 || result[31] != 0x20 {
		return "", err
	}
	res := result[64 : 64+int(result[63])]
	// fmt.Printf("byte: %+v\n", result)
	return string(res), nil
}

func (c *EthClient) GetNonce(addr, param string) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	method := "eth_getTransactionCount"
	var strNonce string
	err := c.RpcClient.CallContext(ctx, &strNonce, method, addr, param)
	if err != nil {
		return 0, err
	}
	nonce := utils.HexTobigInt(strNonce)
	return nonce.Uint64(), nil
}

func (c *EthClient) ChainID() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	id, err := c.Client.ChainID(ctx)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (c *EthClient) GetBlockHeight() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	method := "eth_blockNumber"
	var strHeight string
	err := c.RpcClient.CallContext(ctx, &strHeight, method)
	if err != nil {
		return 0, err
	}
	height := utils.HexTobigInt(strHeight)
	return height.Int64(), nil
}

func (c *EthClient) SuggestGasPrice() *big.Int {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	gasPrice, err := c.Client.SuggestGasPrice(ctx)
	if err != nil {
		return big.NewInt(0)
	}

	return gasPrice
}

func (c *EthClient) GetContractInfoByFunc(contractAddr, params string, args ...interface{}) (interface{}, error) {
	var res []byte
	// ERC20 通用ABI调用
	if len(contractAddr) != 42 {
		return res, fmt.Errorf("invalid contract address length %s", contractAddr)
	}
	// 获取合约信息，转化参数类型
	contractAbi, argsNew, err := GetAbiAndArgs(contractAddr, AbiMap[contractAddr], params, args)
	if err != nil {
		return nil, fmt.Errorf("GetAbiAndArgs error: %+v\n", err)
	}
	data, err := contractAbi.Pack(params, argsNew...)
	if err != nil {
		return res, err
	}
	result, err := c.EvmCall(contractAddr, contractAddr, data)
	if err != nil {
		return res, err
	}

	return contractAbi.Unpack(params, result)
}
