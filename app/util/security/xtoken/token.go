package xtoken

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/copier"
	manager_config "github.com/woshilaixuex/csd_chat_backend/app/manager/config"
	"github.com/woshilaixuex/csd_chat_backend/app/util"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: token生成器
 * @Date: 2025-02-19 20:25
 */
type TokenOption struct {
	SecretKey    string
	AccessExpire int64 // 单位：/s
}

var tokenOption *TokenOption

func DefaultToken() {
	tokenOption.SecretKey = "token"
	tokenOption.AccessExpire = 86400 // 一天
}
func InitJwtToken() {
	tokenOption = new(TokenOption)
	config := manager_config.ConfigsMap[manager_config.TokenConfigName]
	if config == nil {
		DefaultToken()
		return
	}
	configEntity, ok := config.(*manager_config.TokenConfig)
	if !ok {
		panic("")
	}
	err := copier.Copy(tokenOption, configEntity)
	if err != nil {
		panic(err)
	}
}
func GetJwtToken(userId int64) (string, error) {
	if tokenOption == nil {
		InitJwtToken()
	}
	claims := make(jwt.MapClaims)
	now := time.Now().Unix()
	claims["exp"] = now + tokenOption.AccessExpire
	claims["iat"] = now
	claims[util.CtxKeyJwtUserId] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(tokenOption.SecretKey))
}
func ParseUserTokem(token *jwt.Token) (interface{}, error) {
	return []byte(tokenOption.SecretKey), nil
}
