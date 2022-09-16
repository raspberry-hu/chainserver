package config

import (
	"os"

	"github.com/sea-project/go-logger"
	"github.com/spf13/viper"
)

var Conf *Config
var Rpc string
var ChainTypeMap = make(map[string]interface{})
var NftAddrMap = make(map[string]interface{})
var TopicList []string

// Initialize
func Initialize(log_path string) {
	// 加载配置文件
	// 缺省配置
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("./conf/")

	// 自定义配置
	_, err := os.Stat("./conf/my.yaml")
	if err == nil || os.IsExist(err) {
		viper.SetConfigName("my")
		viper.AddConfigPath("./conf/")
		viper.MergeInConfig()
	}

	if err = viper.ReadInConfig(); err != nil {
		logger.Fatal("ReadInConfig error:", err)
	}
	chain_type_list := viper.Get("chain_type_list").([]interface{})
	for _, chain_info := range chain_type_list {
		chain_info_map := chain_info.(map[interface{}]interface{})
		for k, v := range chain_info_map {
			ChainTypeMap[k.(string)] = v
		}
	}
	NftAddrMap = viper.GetStringMap("nft_addr")
	TopicList = viper.GetStringSlice("topic_list")
	Conf = &Config{
		LogCfg: viper.GetString("log"),
		Console: &Console{
			Name:    viper.GetString("console.name"),
			Version: viper.GetString("console.version"),
			Port:    viper.GetString("console.port"),
		},
		Mysql: &Mysql{
			MaxOpenConns: viper.GetInt("mysql.max_open_conns"),
			MaxIdleConns: viper.GetInt("mysql.max_idle_conns"),
			SourceName:   viper.GetString("mysql.source_name"),
		},
		Pgsql: &Pgsql{
			Host:     viper.GetString("pgsql.pghost"),
			User:     viper.GetString("pgsql.pg_user"),
			PassWord: viper.GetString("pgsql.pg_pwd"),
			DbName:   viper.GetString("pgsql.dbname"),
			Port:     viper.GetInt("pgsql.dbport"),
		},
		Jwt: &JWT{
			SigningKey: viper.GetString("jwt.signing-key"),
		},
	}
	Rpc = viper.GetStringSlice("console.nodelist")[0]
	// 初始化日志配置
	logger.SetLogger(log_path)
	logger.Info("Successful configuration load")
}
