package models

type Server struct {
	Service      ServiceConf `mapstructure:"service_conf"  json:"service_conf"  yaml:"service_conf"`  // 服务配置
	Zap          Zap         `mapstructure:"zap"           json:"zap"           yaml:"zap"`           // 日志配置
	UtxoFilepath string      `mapstructure:"utxo_filepath" json:"utxo_filepath" yaml:"utxo_filepath"` // utxo存放路径
	//Mysql       db.Mysql              `mapstructure:"mysql"        json:"mysql"        yaml:"mysql"`        // 数据库配置
	//Https       Https                 `mapstructure:"https"        json:"https"        yaml:"https"`        // 网络配置
	//ElasticConf elastic.ElasticConfig `mapstructure:"elastic_conf" json:"elastic_conf" yaml:"elastic_conf"` // elastic配置
	//ChainNode   ChainNode             `mapstructure:"chain_node"   json:"chain_node"   yaml:"chain_node"`   // 链节点配置
	//CronTasks   CronTasks             `mapstructure:"cron_tasks"   json:"cron_tasks"   yaml:"cron_tasks"`   // 定时任务
}

type ServiceConf struct {
	ServiceAddr           string `mapstructure:"service_addr"            json:"service_addr"            yaml:"service_addr"`
	ServiceFeeAddress     string `mapstructure:"service_fee_address"     json:"service_fee_address"     yaml:"service_fee_address"`
	ServiceFee            int64  `mapstructure:"service_fee"             json:"service_fee"             yaml:"service_fee"`
	ServicePrometheusAddr string `mapstructure:"service_prometheus_addr" json:"service_prometheus_addr" yaml:"service_prometheus_addr"`
}

type Zap struct {
	Level         string `mapstructure:"level"          json:"level"         yaml:"level"`          // 日志级别
	Format        string `mapstructure:"format"         json:"format"        yaml:"format"`         // 输出方式
	Prefix        string `mapstructure:"prefix"         json:"prefix"        yaml:"prefix"`         // 前缀
	Director      string `mapstructure:"director"       json:"director"      yaml:"director"`       // 目录
	LinkName      string `mapstructure:"link-name"      json:"linkName"      yaml:"link-name"`      // 文件名
	ShowLine      bool   `mapstructure:"show-line"      json:"showLine"      yaml:"showLine"`       // 是否展示行号
	EncodeLevel   string `mapstructure:"encode-level"   json:"encodeLevel"   yaml:"encode-level"`   // 日志编码类型
	StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktraceKey" yaml:"stacktrace-key"` // 堆栈跟踪
	LogInConsole  bool   `mapstructure:"log-in-console" json:"logInConsole"  yaml:"log-in-console"` // 是否在工作台输出日志
}
