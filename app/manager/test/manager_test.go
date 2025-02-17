package test

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/woshilaixuex/csd_chat_backend/app/manager/config"
	"github.com/woshilaixuex/csd_chat_backend/app/util/model/manager"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2025-02-17 19:33
 */

// @Config
// 测试初始化 Viper 配置
func setupViper() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

// 测试 Bind 方法是否能正确绑定配置
func TestOrmConfigBind(t *testing.T) {
	// 设置测试数据
	err := setupViper()
	t.Log(err)
	// 创建 OrmConfig 实例
	ormConfig := config.NewOrmConfig()

	// 调用 Bind 方法进行绑定
	err = ormConfig.Bind()

	// 断言 Bind 方法没有返回错误
	assert.NoError(t, err)

	// 断言绑定后的配置值是否正确
	assert.Equal(t, "my_database", ormConfig.Name)
	assert.Equal(t, "mysql", ormConfig.Type)
	t.Log(ormConfig)
}

// @Model
func TestUserManager(t *testing.T) {

	err := manager.InitEngine()
	assert.NoError(t, err)

	user, err := manager.GetUserByID(1)
	assert.NoError(t, err)
	t.Log(user)
}

// @Redis
