package client

import (
	"github.com/ethereum/go-ethereum/common"
)

func EthAddressChange(addr string) common.Address {
	return common.HexToAddress(addr)
}

func HexToHash(hash string) common.Hash {
	return common.HexToHash(hash)
}
