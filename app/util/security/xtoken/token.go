package xtoken

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/copier"
	"github.com/woshilaixuex/csd_chat_backend/app/util"
	"github.com/woshilaixuex/csd_chat_backend/app/util/xconfig"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: token生成器，使用记得初始化一下
 * @Date: 2025-02-19 20:25
 */
type TokenOption struct {
	SecretKey    string
	AccessExpire int64 // 单位：/s
}
type JwtUserToken struct {
	UserId  uint64
	IatTime time.Time
	ExpTime time.Time
}

var tokenOption *TokenOption

func DefaultToken() {
	tokenOption.SecretKey = "token"
	tokenOption.AccessExpire = 86400 // 一天
}
func InitJwtToken() {
	tokenOption = new(TokenOption)
	config := xconfig.ConfigsMap[xconfig.TokenConfigName]
	if config == nil {
		DefaultToken()
		return
	}
	configEntity, ok := config.(*xconfig.TokenConfig)
	if !ok {
		panic("")
	}
	if configEntity.AccessExpire == 0 {
		DefaultToken()
		return
	}
	err := copier.Copy(tokenOption, configEntity)
	if err != nil {
		panic(err)
	}
}
func GetJwtToken(userId uint64) (string, error) {
	if tokenOption == nil {
		InitJwtToken()
	}
	claims := make(jwt.MapClaims)
	now := time.Now().UTC()
	expTime := now.Add(time.Duration(tokenOption.AccessExpire) * time.Second)
	claims["exp"] = expTime.Unix()
	claims["iat"] = now.Unix()
	claims[util.CtxKeyJwtUserId] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	fmt.Println("Current time:", time.Now().UTC())
	fmt.Println("Token expires at:", expTime)

	return token.SignedString([]byte(tokenOption.SecretKey))
}

func ParseUserTokem(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(tokenOption.SecretKey), nil
}

func ParseJwtToken(tokenString string) (uint64, error) {
	if tokenOption == nil {
		InitJwtToken()
	}
	token, err := jwt.Parse(tokenString, ParseUserTokem)

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userId, ok := claims[util.CtxKeyJwtUserId].(float64); ok {
			return uint64(userId), nil
		}
	}

	return 0, fmt.Errorf("invalid token")
}
