package xredis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	manager_config "github.com/woshilaixuex/csd_chat_backend/app/manager/config"
	"github.com/woshilaixuex/csd_chat_backend/app/util"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: redis客户端
 * @Date: 2025-02-18 18:50
 */

type RedisLink struct {
	redisDb *redis.Client
	ctx     context.Context
}

var redisDb *redis.Client

func InitRedisCli() error {
	var err error
	config := manager_config.ConfigsMap[manager_config.RedisConfigName]
	configEntity, ok := config.(*manager_config.RedisConfig)
	if !ok {
		return fmt.Errorf("failed to cast config to RedisConfig")
	}
	redisDb = redis.NewClient(&redis.Options{
		Addr:     configEntity.Address,
		Password: configEntity.Password,
		DB:       configEntity.DB,
	})
	_, err = redisDb.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	return nil
}

// 获取全局CsdID
func GetNewGlobalCsdID() (uint64, error) {
	var err error
	cmd := redisDb.Incr(context.Background(), util.RedisGlobalCsdIdKey)
	if cmd.Err() != nil {
		return 0, err
	}
	id, err := cmd.Uint64()
	if cmd.Err() != nil {
		return 0, err
	}
	return id, err
}
func Set(key string, value interface{}) error {
	var err error
	cmd := redisDb.Set(context.Background(), key, value, -1)
	str, err := cmd.Result()
	log.Printf("redis set key %s: %s", key, str)
	if err != nil {
		return err
	}
	return nil
}

// 获取redis服务器的时间
func GetTimeFromRedis() (time.Time, error) {
	var err error
	cmd := redisDb.Time(context.Background())
	if err = cmd.Err(); err != nil {
		return time.Time{}, err
	}

	time := cmd.Val()

	return time, err
}
