package xconfig

import (
	"github.com/spf13/viper"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2025-02-19 20:30
 */
var TokenConfigName = "token"

type TokenConfig struct {
	SecretKey    string
	AccessExpire int64 // 单位：/s
}

func NewTokenConfig() *TokenConfig {
	return new(TokenConfig)
}
func (c *TokenConfig) Bind() error {
	c.SecretKey = viper.GetString("token.secretkey")
	c.AccessExpire = viper.GetInt64("token.accessexpire")
	return nil
}
func (c *TokenConfig) GetConfigName() string {
	return TokenConfigName
}
