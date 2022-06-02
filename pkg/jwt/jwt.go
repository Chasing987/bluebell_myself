package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

//const TokenExpireDuration = time.Hour * 2

var mySecret = []byte("夏天夏天悄悄过去")

// MyClaims 自定义声明结构体并内嵌 jwt.StandardClaims
// jwt 包自带的 jwt.StandardClaims 只包含了官方字段
// 我们这里需要额外记录一个 username 字段，所以要自定义结构体
// 如果要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenToken 生成JWT
func GenToken(userID int64, username string) (string, error) {
	//创建一个我们自己的声明数据
	c := MyClaims{
		userID,
		username, //自定义字段
		jwt.StandardClaims{
			//ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), //过期时间
			ExpiresAt: time.Now().Add(
				time.Duration(viper.GetInt("auth.jwt_expire")) * time.Hour).Unix(), //过期时间
			Issuer: "bluebell", //签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的 secret 签名并获得完整的编码的字符串 token
	return token.SignedString(mySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	//解析token
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc,
		func(token *jwt.Token) (interface{}, error) {
			return mySecret, nil
		})
	if err != nil {
		return nil, err
	}
	if token.Valid { //检验 token
		return mc, nil
	}
	return nil, errors.New("invalid token")
}
