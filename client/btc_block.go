package client

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"log"
)

// 查询最新区块高度
func (c *BtcClient) GetBlockHeight() (int64, error) {
	return c.Client.GetBlockCount()
}

// 根据块高查HASH
func (c *BtcClient) GetBlockHashByHeight(height int64) (string, error) {
	hash, err := c.Client.GetBlockHash(height)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return hash.String(), err
}

// 根据块高查询数据
func (c *BtcClient) GetBlockInfoByHeight(height int64) (interface{}, error) {
	hash, err := c.Client.GetBlockHash(height)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return c.Client.GetBlock(hash)
}

// 根据块HASH查询数据
func (c *BtcClient) GetBlockInfoByHash(hash string) (interface{}, error) {
	// hash处理
	h, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return c.Client.GetBlock(h)
}
