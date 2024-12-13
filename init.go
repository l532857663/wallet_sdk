package wallet_sdk

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"wallet_sdk/global"
	"wallet_sdk/models"
	"wallet_sdk/utils/dir"
	"wallet_sdk/utils/logutils"
)

func MustLoad(confPath string) {
	global.CONFIG = readConfig(confPath)

	// 初始化zap日志库
	global.LOG = logutils.Log("", global.CONFIG.Zap)

	// 初始化节点类型
	global.ChainName = BTC_RegTest

	//// 初始化elastic数据库
	//elastic.InitElasticInfo(global.CONFIG.ElasticConf)

	// 初始化本地存放UTXO文件夹
	initUTXOPath()
}

func readConfig(confPath string) *models.Server {
	var config *models.Server
	data, err := os.ReadFile(confPath)
	if err != nil {
		log.Fatalf("Read conf file[%s] error: %v", confPath, err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Yaml unmarshal error: %v", err)
	}
	return config
}

func initUTXOPath() {
	// 初始化
	global.UtxoBlockHeightPath = global.CONFIG.UtxoFilepath + "/block"
	global.UtxoSpendPath = global.CONFIG.UtxoFilepath + "/spend"
	global.UtxoUnSpendPath = global.CONFIG.UtxoFilepath + "/unSpend"
	pathList := []string{
		global.CONFIG.UtxoFilepath,
		global.UtxoSpendPath,
		global.UtxoUnSpendPath,
	}
	if err := dir.CreateDir(pathList...); err != nil {
		logutils.LogErrorf(global.LOG, "Error creating UTXO file:%v", err)
		return
	}
}

func Shutdown() {
}
