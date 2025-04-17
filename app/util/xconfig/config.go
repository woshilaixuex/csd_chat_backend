package xconfig

import (
	"github.com/spf13/viper"
	_ "github.com/woshilaixuex/csd_chat_backend/app/util/xlog"
	"golang.org/x/exp/slog"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 配置集合
 * @Date: 2025-02-17 07:19
 */

// Config 配置接口
type Config interface {
	GetConfigName() string
	Bind() error
}

// Configs 配置集合
var ConfigsMap = make(map[string]Config)

// addConfigs 添加配置到配置集合
func AddConfigs(configs ...Config) {
	for _, config := range configs {
		err := config.Bind()
		if err != nil {
			slog.Error("config", "err", err.Error())
			panic(err)
		}
		ConfigsMap[config.GetConfigName()] = config
	}
}

// 初始化 Viper 配置
func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../")
	err := viper.ReadInConfig()
	if err != nil {
		slog.Error("config", "err", err.Error())
	} else {
		slog.Info("Config file loaded successfully")
	}

	AddConfigs(NewOrmConfig(),
		NewRedisConfig(),
		NewTokenConfig(),
		NewEtcdConfig())
}
