package xconfig

import "github.com/spf13/viper"

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 数据库映射配置
 * @Date: 2025-02-16 21:59
 */

var OrmConifgName = "orm"

type OrmConfig struct {
	Name string
	Type string
	Url  string
}

func NewOrmConfig() *OrmConfig {
	return new(OrmConfig)
}
func (c *OrmConfig) Bind() error {
	c.Name = viper.GetString("database.name")
	c.Type = viper.GetString("database.type")
	c.Url = viper.GetString("database.url")
	return nil
}

func (c *OrmConfig) GetConfigName() string {
	return OrmConifgName
}
