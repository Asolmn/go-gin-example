module github.com/Asolmn/go-gin-example

go 1.19

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fvbock/endless v0.0.0-20170109170031-447134032cb6
	github.com/gin-gonic/gin v1.9.0
	github.com/go-ini/ini v1.67.0
	github.com/gomodule/redigo v1.8.9
	github.com/robfig/cron v1.2.0
	github.com/swaggo/files v1.0.1
	github.com/swaggo/gin-swagger v1.6.0
	github.com/swaggo/swag v1.8.12
	github.com/unknwon/com v1.0.1
	gorm.io/driver/mysql v1.4.7
	gorm.io/gorm v1.24.6
)

// validator接口里会使用到表单验证
require github.com/astaxie/beego v1.10.1

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/boombuler/barcode v1.0.1 // indirect
	github.com/bytedance/sonic v1.8.7 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/spec v0.20.8 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.12.0 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/leodido/go-urn v1.2.2 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/pelletier/go-toml/v2 v2.0.7 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.3 // indirect
	github.com/tealeg/xlsx v1.0.5 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	github.com/xuri/efp v0.0.0-20220603152613-6918739fd470 // indirect
	github.com/xuri/excelize/v2 v2.7.1 // indirect
	github.com/xuri/nfp v0.0.0-20220409054826-5e722a1d9e22 // indirect
	golang.org/x/arch v0.3.0 // indirect
	golang.org/x/crypto v0.8.0 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	golang.org/x/tools v0.7.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/ASOLMN/go-gin-example/conf => ./conf

	github.com/ASOLMN/go-gin-example/middleware => ./middleware

	github.com/ASOLMN/go-gin-example/models => ./models

	github.com/ASOLMN/go-gin-example/pkg/setting => ./pkg/setting

	github.com/ASOLMN/go-gin-example/routers => ./routers
	github.com/Asolmn/go-gin-example/docs => ./docs
	github.com/Asolmn/go-gin-example/middleware/jwt => ./middleware/jwt
	github.com/Asolmn/go-gin-example/pkg/app => ./pkg/app
	github.com/Asolmn/go-gin-example/pkg/e => ./pkg/e

	github.com/Asolmn/go-gin-example/pkg/export => ./pkg/export
	github.com/Asolmn/go-gin-example/pkg/file => ./pkg/file
	github.com/Asolmn/go-gin-example/pkg/gredis => ./pkg/gredis
	github.com/Asolmn/go-gin-example/pkg/upload => ./pkg/upload
	github.com/Asolmn/go-gin-example/pkg/util => ./pkg/util
	github.com/Asolmn/go-gin-example/routers/api => ./routers/api
	github.com/Asolmn/go-gin-example/service/article_service => ./service/article_service

	github.com/Asolmn/go-gin-example/service/cache_service => ./service/cache_service
)
