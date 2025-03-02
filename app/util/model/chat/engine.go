package chat

import (
	chat_config "github.com/woshilaixuex/csd_chat_backend/app/im/config"
	"xorm.io/xorm"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2025-02-22 21:14
 */
var engine *xorm.Engine

func InitEngine() error {
	var err error
	config := chat_config.ConfigsMap[chat_config.OrmConifgName]
	configEntity, ok := config.(*chat_config.OrmConfig)
	if !ok {
		return err
	}
	engine, err = xorm.NewEngine(configEntity.Type, configEntity.Url)
	if err != nil {
		return err
	}
	err = engine.Ping()
	return err
}
