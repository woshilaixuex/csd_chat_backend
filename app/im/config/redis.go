package config

import (
	"time"

	"github.com/spf13/viper"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: redis配置
 * @Date: 2025-02-17 21:47
 */
var RedisConfigName = "redis"

type RedisConfig struct {
	Name     string        // 可选字段，Redis 配置的名称（通常是用于多个 Redis 实例区分）
	Address  string        // Redis 服务器地址，如 "localhost:6379"
	Password string        // Redis 密码，若没有密码可为空
	DB       int           // Redis 数据库索引，默认是 0
	PoolSize int           // Redis 连接池的大小，决定并发请求的数量
	MinIdle  int           // 最小空闲连接数
	MaxIdle  int           // 最大空闲连接数
	Timeout  time.Duration // Redis 连接的超时时间
}

func NewRedisConfig() *RedisConfig {
	return new(RedisConfig)
}

func (c *RedisConfig) Bind() error {
	c.Name = viper.GetString("redis.name")
	c.Address = viper.GetString("redis.address")
	c.Password = viper.GetString("redis.password")
	c.DB = viper.GetInt("redis.db")
	c.PoolSize = viper.GetInt("redis.pool_size")
	c.MinIdle = viper.GetInt("redis.min_idle")
	c.MaxIdle = viper.GetInt("redis.max_idle")
	c.Timeout = viper.GetDuration("redis.timeout")
	return nil
}

func (c *RedisConfig) GetConfigName() string {
	return RedisConfigName
}
