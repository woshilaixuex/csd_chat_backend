package xredis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/woshilaixuex/csd_chat_backend/app/util"
	"github.com/woshilaixuex/csd_chat_backend/app/util/xconfig"
	"github.com/woshilaixuex/csd_chat_backend/app/util/xerr"
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
	initErr := xerr.RedisBackGroundError
	config := xconfig.ConfigsMap[xconfig.RedisConfigName]
	configEntity, ok := config.(*xconfig.RedisConfig)
	if !ok {
		return initErr.Wrap(fmt.Errorf("failed to cast config to RedisConfig"))
	}
	redisDb = redis.NewClient(&redis.Options{
		Addr:     configEntity.Address,
		Password: configEntity.Password,
		DB:       configEntity.DB,
	})
	_, err := redisDb.Ping(context.Background()).Result()
	if err != nil {
		return initErr.Wrap(err)
	}
	return initErr.Submit()
}

// 获取全局CsdID
func GetNewGlobalCsdID() (uint64, error) {
	backErr := xerr.RedisBackGroundError
	cmd := redisDb.Incr(context.Background(), util.RedisGlobalCsdIdKey)
	if cmd.Err() != nil {
		return 0, backErr.Wrap(cmd.Err())
	}
	id, err := cmd.Uint64()
	if err != nil {
		return 0, backErr.Wrap(cmd.Err())
	}
	return id, backErr.Submit()
}

func Set(key string, value interface{}) error {
	backErr := xerr.RedisBackGroundError
	cmd := redisDb.Set(context.Background(), key, value, -1)
	str, err := cmd.Result()
	log.Printf("redis set key %s: %s", key, str)
	if err != nil {
		return backErr.Wrap(err)
	}
	return backErr.Submit()
}

// 获取redis服务器的时间
func GetTimeFromRedis() (time.Time, error) {
	backErr := xerr.RedisBackGroundError
	cmd := redisDb.Time(context.Background())
	if err := cmd.Err(); err != nil {
		return time.Time{}, backErr.Wrap(err)
	}

	time := cmd.Val()

	return time, backErr.Submit()
}
