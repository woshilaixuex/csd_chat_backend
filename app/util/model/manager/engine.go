package manager

import (
	_ "github.com/go-sql-driver/mysql"
	manager_config "github.com/woshilaixuex/csd_chat_backend/app/manager/config"
	"xorm.io/xorm"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 用户管理数据库orm
 * @Date: 2025-02-17 19:21
 */
var engine *xorm.Engine

func InitEngine() error {
	var err error
	config := manager_config.ConfigsMap[manager_config.OrmConifgName]
	configEntity, ok := config.(*manager_config.OrmConfig)
	if !ok {
		return err
	}
	engine, err = xorm.NewEngine(configEntity.Type, configEntity.Url)
	if err != nil {
		return err
	}
	return err
}
