package test

import (
	"testing"

	"github.com/jinzhu/copier"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/woshilaixuex/csd_chat_backend/app/manager/config"
	"github.com/woshilaixuex/csd_chat_backend/app/manager/internal/user"
	"github.com/woshilaixuex/csd_chat_backend/app/util/model/manager"
	"github.com/woshilaixuex/csd_chat_backend/app/util/security/encryption"
	"github.com/woshilaixuex/csd_chat_backend/app/util/security/xtoken"
	"github.com/woshilaixuex/csd_chat_backend/app/util/xredis"
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

	_, user, err := manager.GetUserByID(0)
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

	userId := uint64(123456)
	token, err := xtoken.GetJwtToken(userId)
	t.Log(token)
	assert.NoError(t, err)

	id, err := xtoken.ParseJwtToken(token)
	assert.NoError(t, err)
	t.Log(id)
}

// @Encrypt
func TestEncrypt(t *testing.T) {
	entity := &user.RegisterEntity{
		UserName:    "nihao",
		StudentID:   "11",
		RealName:    "小明",
		PhoneNumber: "1233",
		Email:       "11@qq.com",
		Password:    "123",
	}
	salt, hash, err := encryption.EncryptPassword(entity.Password)
	if err != nil {
		return
	}
	entity.Salt = salt
	entity.Password = hash
	user := new(manager.UserManager)
	copier.Copy(user, entity)
	t.Log(user)
	isTrue := encryption.VerifyPassword(user.HashPassword, "123", entity.Salt)
	t.Log(isTrue)
}

// @Serivce:user
func TestSerivceUser(t *testing.T) {
	manager.InitEngine()
	xredis.InitRedisCli()
	entity := &user.RegisterEntity{
		UserName:    "nihao3",
		StudentID:   "113",
		RealName:    "小明",
		PhoneNumber: "1233",
		Email:       "131@qq.com",
		Password:    "123",
	}
	token, err := user.Register(entity)
	assert.NoError(t, err)
	t.Log(token)
}
