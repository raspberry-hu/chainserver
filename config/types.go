package config

// yaml配置结构体
type Config struct {
	LogCfg   string   `json:"log"`
	Console  *Console `json:"console"`
	Mysql    *Mysql   `json:"mysql"`
	Pgsql    *Pgsql   `json:"pgsql"`
	NodeList []string `json:"nodeList"`
	Jwt      *JWT     `json:"jwt"`
}

// yaml配置结构体
type Console struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Port    string `json:"port"`
}

type Mysql struct {
	MaxOpenConns int    `json:"max_open_conns"`
	MaxIdleConns int    `json:"max_idle_conns"`
	SourceName   string `json:"source_name"`
}

type Pgsql struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	PassWord string `json:"password"`
	DbName   string `json:"dbname"`
	Port     int    `json:"port"`
}

type JWT struct {
	SigningKey string `mapstructure:"signing-key" json:"signingKey"`
}

/**
与蚂蚁链相关的结构体
*/
type TransHash struct {
	TxHash string `json:"txHash"`
}
type Output struct {
	OutRes      []interface{} `json:"outRes"`
	Transaction TransHash     `json:"transaction"`
}
type TxTime struct {
	Time interface{} `json:"time"`
}
