package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func EncodeStringByUtxoInfo(txHash string, vout uint32, amount int64) string {
	return fmt.Sprintf("%s:%d %d", txHash, vout, amount)
}

func DecodeUtxoInfoByString(val string) (string, uint32, int64) {
	var txHash string
	var vout uint32
	var amount int64
	utxoInfo := strings.Split(val, " ")
	if len(utxoInfo) == 0 {
		return txHash, vout, amount
	}
	utxo := strings.Split(utxoInfo[0], ":")
	if len(utxo) != 2 {
		return txHash, vout, amount
	}
	txHash = utxo[0]
	v, _ := strconv.ParseUint(utxo[1], 0, 0)
	amount, _ = strconv.ParseInt(utxoInfo[1], 0, 0)
	return txHash, uint32(v), amount
}

// filterData 用于根据查询字符串过滤数据
func FilterData(data []string, query string) []string {
	var result []string
	for _, item := range data {
		if strings.Contains(strings.ToLower(item), strings.ToLower(query)) {
			result = append(result, item)
		}
	}
	return result
}
