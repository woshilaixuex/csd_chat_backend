package xetcd

import (
	"context"
	"fmt"
	"sync"

	"github.com/woshilaixuex/csd_chat_backend/app/util/xerr"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdLink struct {
	cli         *clientv3.Client
	renewalTime int64
	mu          sync.RWMutex                // 读写锁，用于并发安全
	leases      map[string]clientv3.LeaseID // 管理已注册服务的租约
}

// etcd:默认信息
const (
	defaultRenewalTime int64 = 20
)

// etcd:全局
var (
	etcdOnce sync.Once
	etcdCon  *EtcdLink
)

func InitEtcdLink() error {
	initErr := xerr.EtcdBackGroundError
}

func NewEtcdLink(config clientv3.Config) (*EtcdLink, error) {
	backErr := xerr.EtcdBackGroundError
	etcdOnce.Do(func() {
		if len(config.Endpoints) == 0 {
			backErr.Wrap(xerr.EtcdErrInvalidConfig)
			return
		}

		cli, err := clientv3.New(config)
		if err != nil {
			backErr.Wrap(xerr.EtcdErrNotInitialized.Err)
			return
		}

		etcdCon = &EtcdLink{
			cli:         cli,
			renewalTime: defaultRenewalTime,
			leases:      make(map[string]clientv3.LeaseID),
		}
	})
	return etcdCon, backErr.Submit()
}

func (e *EtcdLink) Close() error {
	err := xerr.EtcdBackGroundError
	if e.cli == nil {
		return err.Wrap(xerr.EtcdErrNotInitialized)
	}
	return e.cli.Close()
}

func (e *EtcdLink) SetRenewalTime(seconds int64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.renewalTime = seconds
}

// RegisterService 注册服务（含自动续约）
func (e *EtcdLink) RegisterService(serviceKey, addr string) error {
	backErr := xerr.EtcdBackGroundError
	if e.cli == nil {
		return backErr.Wrap(xerr.EtcdErrNotInitialized)
	}

	// 1. 生成租约
	lease := clientv3.NewLease(e.cli)
	grantResp, err := lease.Grant(context.TODO(), e.renewalTime)
	if err != nil {
		return backErr.Wrap(fmt.Errorf("create lease failed: %w", err))
	}

	// 2. 写入键值（绑定租约）
	key := fmt.Sprintf("/services/%s/%s", serviceKey, addr)
	_, err = e.cli.Put(context.TODO(), key, addr, clientv3.WithLease(grantResp.ID))
	if err != nil {
		return backErr.Wrap(fmt.Errorf("put key failed: %w", err))
	}

	// 3. 启动续约协程
	keepAlive, err := lease.KeepAlive(context.TODO(), grantResp.ID)
	if err != nil {
		return backErr.Wrap(fmt.Errorf("start keepalive failed: %w", err))
	}

	// 4. 管理租约状态
	e.mu.Lock()
	e.leases[key] = grantResp.ID
	e.mu.Unlock()

	// 5. 后台续约
	go func() {
		for range keepAlive {
			// 自动续约，通道关闭时触发重试逻辑
		}
	}()

	return nil
}

func (e *EtcdLink) UnregisterService(serviceKey, addr string) error {
	backErr := xerr.EtcdBackGroundError
	key := fmt.Sprintf("/services/%s/%s", serviceKey, addr)

	e.mu.Lock()
	defer e.mu.Unlock()

	if leaseID, ok := e.leases[key]; ok {
		if _, err := e.cli.Revoke(context.TODO(), leaseID); err != nil {
			return backErr.Wrap(err)
		}
		delete(e.leases, key)
	}
	_, err := e.cli.Delete(context.TODO(), key)
	return backErr.Wrap(err).Submit()
}
