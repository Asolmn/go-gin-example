package util

import (
	"github.com/Asolmn/go-gin-example/pkg/setting"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtSecret = []byte(setting.AppSetting.JwtSecret)

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// 生成token
func GenerateToken(username, password string) (string, error) {

	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		EncodeMD5(username),
		EncodeMD5(password),
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), // 生效时间
			Issuer:    "gin-blog",        // 颁布者
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	//fmt.Println(err)
	return token, err
}

// 验证token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		// 类型断言，如果是属于Claims类型, 返回动态值和ok
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid { // Valid验证基于时间的声明
			return claims, nil
		}
	}
	return nil, err
}
