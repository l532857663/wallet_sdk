package client

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// 根据合约方法处理参数类型
func GetAbiAndArgs(abiContent, params string, args []interface{}) (abi.ABI, []interface{}, error) {
	// Abi转化
	contractAbi, err := StringToAbi(abiContent)
	if err != nil {
		return contractAbi, nil, err
	}
	method, exist := contractAbi.Methods[params]
	fmt.Printf("GetAbiAndArgs method: %+v\n, exist: %+v\n", method, exist)
	var argsNew []interface{}
	abiParam := method.Inputs
	for i, v := range abiParam {
		arg := ChangeArgType(args[i], v.Type.T)
		_, ok := arg.(error)
		if arg == nil || ok {
			continue
		}
		argsNew = append(argsNew, arg)
	}
	// 检查参数数量是否匹配
	if len(argsNew) != len(abiParam) {
		err := fmt.Errorf("the args len not enough")
		return contractAbi, nil, err
	}
	return contractAbi, argsNew, nil
}

func ChangeArgType(arg interface{}, argType byte) interface{} {
	argStr := arg.(string)
	switch argType {
	case abi.AddressTy:
		addr := EthAddressChange(argStr)
		if addr.String() != argStr {
			return nil
		}
		return addr
	case abi.UintTy:
		val, ok := big.NewInt(0).SetString(argStr, 10)
		if !ok {
			return nil
		}
		return val
	case abi.SliceTy:
		var argSlice []common.Address
		err := json.Unmarshal([]byte(argStr), &argSlice)
		if err != nil {
			return err
		}
		return argSlice
	default:
		fmt.Printf("ChangeArgType not case: %+v\n", argType)
	}
	return nil
}

func EthAddressChange(addr string) common.Address {
	return common.HexToAddress(addr)
}

func StringToAbi(abiContent string) (abi.ABI, error) {
	return abi.JSON(strings.NewReader(abiContent))
}

var (
	Erc20Abi, _ = abi.JSON(strings.NewReader(`[
  {
    "constant": true,
    "inputs": [],
    "name": "name",
    "outputs": [
      {
        "name": "",
        "type": "string"
      }
    ],
    "type": "function"
  },
  {
    "constant": false,
    "inputs": [
      {
        "name": "_from",
        "type": "address"
      },
      {
        "name": "_to",
        "type": "address"
      },
      {
        "name": "_value",
        "type": "uint256"
      }
    ],
    "name": "transferFrom",
    "outputs": [
      {
        "name": "success",
        "type": "bool"
      }
    ],
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [],
    "name": "decimals",
    "outputs": [
      {
        "name": "",
        "type": "uint8"
      }
    ],
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [
      {
        "name": "",
        "type": "address"
      }
    ],
    "name": "balanceOf",
    "outputs": [
      {
        "name": "",
        "type": "uint256"
      }
    ],
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [],
    "name": "symbol",
    "outputs": [
      {
        "name": "",
        "type": "string"
      }
    ],
    "type": "function"
  },
  {
    "constant": false,
    "inputs": [
      {
        "name": "_to",
        "type": "address"
      },
      {
        "name": "_value",
        "type": "uint256"
      }
    ],
    "name": "transfer",
    "outputs": [],
    "type": "function"
  },
  {
    "constant": false,
    "inputs": [
      {
        "name": "_spender",
        "type": "address"
      },
      {
        "name": "_value",
        "type": "uint256"
      },
      {
        "name": "_extraData",
        "type": "bytes"
      }
    ],
    "name": "approveAndCall",
    "outputs": [
      {
        "name": "success",
        "type": "bool"
      }
    ],
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [
      {
        "name": "",
        "type": "address"
      },
      {
        "name": "",
        "type": "address"
      }
    ],
    "name": "spentAllowance",
    "outputs": [
      {
        "name": "",
        "type": "uint256"
      }
    ],
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [
      {
        "name": "",
        "type": "address"
      },
      {
        "name": "",
        "type": "address"
      }
    ],
    "name": "allowance",
    "outputs": [
      {
        "name": "",
        "type": "uint256"
      }
    ],
    "type": "function"
  },
  {
    "inputs": [
      {
        "name": "initialSupply",
        "type": "uint256"
      },
      {
        "name": "tokenName",
        "type": "string"
      },
      {
        "name": "decimalUnits",
        "type": "uint8"
      },
      {
        "name": "tokenSymbol",
        "type": "string"
      }
    ],
    "type": "constructor"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "name": "from",
        "type": "address"
      },
      {
        "indexed": true,
        "name": "to",
        "type": "address"
      },
      {
        "indexed": false,
        "name": "value",
        "type": "uint256"
      }
    ],
    "name": "Transfer",
    "type": "event"
  }]`))
)
