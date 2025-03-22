package xetcd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/copier"
	"github.com/woshilaixuex/csd_chat_backend/app/util/xconfig"
	"github.com/woshilaixuex/csd_chat_backend/app/util/xerr"
	clientv3 "go.etcd.io/etcd/client/v3"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 客户端api(借鉴kit源码) 客户端基本功能:get获取去注册信息；watch监听并更新服务端信息；
 * @Date: 2025-03-19 20:41
 */

// Client is a wrapper around the etcd client.
type Client interface {
	// GetEntries queries the given prefix in etcd and returns a slice
	// containing the values of all keys found, recursively, underneath that
	// prefix.
	GetEntries(prefix string) ([]string, error)

	// 监听prefix前缀变化
	WatchPrefix(prefix string, ch chan struct{})

	// 注册对应服务到etcd
	Register(s Service) error

	// 注销etcd的服务
	Deregister(s Service) error

	// 返回为此服务实例创建的租赁id
	LeaseID() int64
}

type client struct {
	cli *clientv3.Client
	ctx context.Context
	kv  clientv3.KV

	// Watcher interface instance, used to leverage Watcher.Close()
	watcher clientv3.Watcher
	// watcher context
	wctx context.Context
	// watcher cancel func
	wcf context.CancelFunc

	// leaseID will be 0 (clientv3.NoLease) if a lease was not created
	leaseID clientv3.LeaseID

	kpAliveCh <-chan *clientv3.LeaseKeepAliveResponse
	// Lease interface instance, used to leverage Lease.Close()
	leaser clientv3.Lease
}

func NewClinet() (Client, error) {
	initErr := xerr.EtcdBackGroundError
	config := xconfig.ConfigsMap[xconfig.EtcdConfigName]
	configEntity, ok := config.(*xconfig.EtcdConfig)
	if !ok {
		return nil, initErr.Wrap(fmt.Errorf("failed to cast config to EtcdConfig"))
	}

	var etcdConfig clientv3.Config
	copier.Copy(etcdConfig, configEntity)
	clientv3.New(etcdConfig)

	return &client{}, initErr.Submit()
}

// 解析bytes value值转成string
func (c *client) GetEntries(prefixKey string) ([]string, error) {
	backErr := xerr.EtcdBackGroundError
	resp, err := c.kv.Get(c.ctx, prefixKey, clientv3.WithPrefix())
	if err != nil {
		return nil, backErr.Wrap(err)
	}
	entries := make([]string, len(resp.Kvs))
	for i, kv := range resp.Kvs {
		entries[i] = string(kv.Value)
	}

	return entries, backErr.Submit()
}

func (c *client) WatchPrefix(prefix string, ch chan struct{}) {
	c.wctx, c.wcf = context.WithCancel(c.ctx)
	c.watcher = clientv3.NewWatcher(c.cli)
	wch := c.watcher.Watch(c.wctx, prefix, clientv3.WithPrefix(), clientv3.WithRev(0))
	ch <- struct{}{}
	for wr := range wch {
		if wr.Canceled {
			return
		}
		ch <- struct{}{}
	}
}

func (c *client) Register(s Service) error {
	backErr := xerr.EtcdBackGroundError

	if s.Key == "" {
		return backErr.Wrap(xerr.EtcdErrNoKey)
	}
	if s.Value == "" {
		return backErr.Wrap(xerr.EtcdErrNoValue)
	}

	// 每次调用时刷新资源
	if c.leaser != nil {
		c.leaser.Close()
	}
	c.leaser = clientv3.NewLease(c.cli)
	if c.watcher != nil {
		c.watcher.Close()
	}
	c.watcher = clientv3.NewWatcher(c.cli)
	if c.kv == nil {
		c.kv = clientv3.NewKV(c.cli)
	}
	if s.TTL == nil {
		s.TTL = NewTTLOption(time.Second*3, time.Second*10)
	}

	// 请求获取租约
	grantResp, err := c.leaser.Grant(c.ctx, int64(s.TTL.ttl.Seconds()))
	if err != nil {
		return backErr.Wrap(err)
	}
	c.leaseID = grantResp.ID
	putResponse, err := c.kv.Put(
		c.ctx,
		s.Key,
		s.Value,
		clientv3.WithLease(c.leaseID),
	)
	log.Printf("debug: register serivce:%s", string(putResponse.PrevKv.Value))
	if err != nil {
		return backErr.Wrap(err)
	}
	// 保活
	c.kpAliveCh, err = c.cli.KeepAlive(c.ctx, c.leaseID)
	if err != nil {
		return backErr.Wrap(err)
	}
	go func() {
		for {
			select {
			case r := <-c.kpAliveCh:
				// avoid dead loop when channel was closed
				if r == nil {
					return
				}
			case <-c.ctx.Done():
				return
			}
		}
	}()

	return backErr.Submit()
}

func (c *client) Deregister(s Service) error {
	backErr := xerr.EtcdBackGroundError
	defer c.close()

	if s.Key == "" {
		return backErr.Wrap(xerr.EtcdErrNoKey)
	}
	if _, err := c.cli.Delete(c.ctx, s.Key, clientv3.WithIgnoreLease()); err != nil {
		return backErr.Wrap(err)
	}

	return backErr.Submit()
}

func (c *client) LeaseID() int64 {
	return int64(c.leaseID)
}

func (c *client) close() {
	if c.leaser != nil {
		c.leaser.Close()
	}
	if c.watcher != nil {
		c.watcher.Close()
	}
	if c.wcf != nil {
		c.wcf()
	}
}
