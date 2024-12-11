package wallet_sdk

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"wallet_sdk/global"
	"wallet_sdk/models"
	"wallet_sdk/utils/logutils"
)

func MustLoad(confPath string) {
	global.CONFIG = readConfig(confPath)

	// 初始化zap日志库
	global.LOG = logutils.Log("", global.CONFIG.Zap)

	//// 初始化elastic数据库
	//elastic.InitElasticInfo(global.CONFIG.ElasticConf)
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

func Shutdown() {
}
