package manager

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/woshilaixuex/csd_chat_backend/app/util/xconfig"
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
	config := xconfig.ConfigsMap[xconfig.OrmConifgName]
	configEntity, ok := config.(*xconfig.OrmConfig)
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
