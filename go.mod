module ChainServer

go 1.14

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/ethereum/go-ethereum v1.10.16
	github.com/gin-gonic/gin v1.7.2
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gohouse/converter v0.0.3
	github.com/google/uuid v1.3.0
	github.com/jinzhu/gorm v1.9.16
	github.com/lib/pq v1.10.3
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/robfig/cron v1.2.0
	github.com/satori/go.uuid v1.2.0
	github.com/sea-project/go-logger v0.0.0-20201130035355-195a8afd0421
	github.com/smartwalle/alipay/v3 v3.1.7
	github.com/smartwalle/xid v1.0.6
	github.com/spf13/viper v1.7.1
	gitlab.alipay-inc.com/antchain/restclient-go-sdk v0.0.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)

replace gitlab.alipay-inc.com/antchain/restclient-go-sdk => ../restclient-go-sdk
