package xetcd

import (
	"fmt"
	"sync"
	"time"

	"github.com/woshilaixuex/csd_chat_backend/app/util/xconfig"
	"github.com/woshilaixuex/csd_chat_backend/app/util/xerr"
	clientv3 "go.etcd.io/etcd/client/v3"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: etcd核心模块
 * @Date: 2025-03-19 20:40
 */

// etcd:默认信息
const (
	defaultRenewalTime time.Duration = 20 * time.Second
)

// etcd:全局
var (
	etcdOnce sync.Once
)

type EtcdClinet struct {
	cli    *clientv3.Client // 默认从一个集群中获取信息
	mu     sync.RWMutex
	leases map[string]clientv3.LeaseID
}

type EtcdSerice struct {
	cli         *clientv3.Client
	renewalTime int64 // 续约时间
}

func newEtcdCon() (*clientv3.Client, error) {
	initErr := xerr.EtcdBackGroundError
	config := xconfig.ConfigsMap[xconfig.EtcdConfigName]
	configEntity, ok := config.(*xconfig.EtcdConfig)
	if !ok {
		return nil, initErr.Wrap(fmt.Errorf("failed to cast config to RedisConfig"))
	}
	if configEntity.Time == 0 {
		configEntity.Time = defaultRenewalTime
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   configEntity.Endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, initErr.Wrap(err)
	}
	return cli, initErr.Submit()
}
