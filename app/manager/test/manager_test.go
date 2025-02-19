package test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/copier"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/woshilaixuex/csd_chat_backend/app/manager/config"
	"github.com/woshilaixuex/csd_chat_backend/app/manager/internal/user"
	"github.com/woshilaixuex/csd_chat_backend/app/util/model/manager"
	xredis "github.com/woshilaixuex/csd_chat_backend/app/util/redis"
	"github.com/woshilaixuex/csd_chat_backend/app/util/security/xtoken"
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
func TestRedisCli(t *testing.T) {

	err := xredis.InitRedisCli()
	assert.NoError(t, err)
	xredis.Set("key_test", 1)
}

// @Copier
func TestCopier(t *testing.T) {
	entity := &user.RegisterEntity{
		UserName: "nihao",
		Password: "123",
	}
	user := new(manager.UserManager)
	copier.Copy(user, entity)
	t.Log(user)
}

// @Jwt
func TestGetJwtToken(t *testing.T) {
	xtoken.InitJwtToken()

	userId := int64(12345)
	token, err := xtoken.GetJwtToken(userId)
	t.Log(token)
	// 断言
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	time.Sleep(time.Second * 2)
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("token"), nil
	})
	assert.NoError(t, err)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	t.Log(claims)
}
