package xetcd

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"time"

	"github.com/woshilaixuex/csd_chat_backend/app/util/xconfig"
	"github.com/woshilaixuex/csd_chat_backend/app/util/xerr"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	GetServerEntries(prefix string) ([]string, error)

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
	// 客户端的心跳保持
	kpAliveCh <-chan *clientv3.LeaseKeepAliveResponse
	leaser    clientv3.Lease
}

// ClientOptions defines options for the etcd client. All values are optional.
// If any duration is not specified, a default of 3 seconds will be used.
type ClientOptions struct {
	Cert          string
	Key           string
	CACert        string
	DialTimeout   time.Duration
	DialKeepAlive time.Duration

	// DialOptions is a list of dial options for the gRPC client (e.g., for interceptors).
	// For example, pass grpc.WithBlock() to block until the underlying connection is up.
	// Without this, Dial returns immediately and connecting the server happens in background.
	DialOptions []grpc.DialOption
}

// 初始化测试连接异常
func Ping(cli *clientv3.Client) error {
	initErr := xerr.EtcdBackGroundError
	_, err := cli.Get(context.Background(), "health-check")
	if err == nil || status.Code(err) == codes.NotFound {
		slog.Info("etcd", "ping", "etcd connected successfully")
	} else {
		return initErr.Wrap(err)
	}
	return initErr.Submit()
}

func NewClinet(ctx context.Context, options ...ClientOptions) (Client, error) {
	// 自定义错误与获取配置信息
	initErr := xerr.EtcdBackGroundError
	config := xconfig.ConfigsMap[xconfig.EtcdConfigName]
	configEntity, ok := config.(*xconfig.EtcdConfig)
	slog.Debug("etcd", "config entity", configEntity)

	if !ok {
		return nil, initErr.Wrap(fmt.Errorf("failed to cast config to EtcdConfig"))
	}

	if len(configEntity.Endpoints)%2 != 1 {
		return nil, initErr.Wrap(xerr.EtcdErrConfigSplitBrain)
	}
	var cli *clientv3.Client

	var err error

	switch configEntity.Method {
	case xconfig.WITHUSERSTL:
		if options == nil {
			return nil, initErr.Wrap(err)
		}
		cli, err = NewClinetWithTLS(ctx, configEntity, options[0])
		if err != nil {
			return nil, initErr.Wrap(err)
		}
	case xconfig.WITHUSERPASSWORD:
		cli, err = NewClinetWithPassword(ctx, configEntity)
		slog.Debug("etcd", "action init with password", err)
		if err != nil {
			return nil, initErr.Wrap(err)
		}
	default:
		slog.Error("etcd", "action init", xerr.EtcdErrUnknownMethod.Error())
		return nil, initErr.Wrap(xerr.EtcdErrUnknownMethod)
	}
	err = Ping(cli)
	if err != nil {
		return nil, initErr.Wrap(err)
	}
	return &client{
		cli: cli,
		ctx: ctx,
		kv:  clientv3.NewKV(cli),
	}, initErr.Submit()
}
func NewClinetWithTLS(ctx context.Context, config *xconfig.EtcdConfig, options ClientOptions) (*clientv3.Client, error) {
	backErr := xerr.EtcdBackGroundError
	var tlscfg *tls.Config
	var err error
	if options.Cert != "" && options.Key != "" {
		tlsInfo := transport.TLSInfo{
			CertFile:      options.Cert,
			KeyFile:       options.Key,
			TrustedCAFile: options.CACert,
		}
		tlscfg, err = tlsInfo.ClientConfig()
		if err != nil {
			return nil, backErr.Wrap(err)
		}
	}
	cli, err := clientv3.New(clientv3.Config{
		Context:           ctx,
		Endpoints:         config.Endpoints,
		DialTimeout:       options.DialTimeout,
		DialKeepAliveTime: options.DialKeepAlive,
		DialOptions:       options.DialOptions,
		TLS:               tlscfg,
		Username:          config.UserName,
		Password:          config.Password,
	})
	if err != nil {
		return nil, backErr.Wrap(err)
	}
	return cli, backErr.Submit()
}
func NewClinetWithPassword(ctx context.Context, config *xconfig.EtcdConfig) (*clientv3.Client, error) {
	backErr := xerr.EtcdBackGroundError
	cli, err := clientv3.New(clientv3.Config{
		Context:   ctx,
		Endpoints: config.Endpoints,
		Username:  config.UserName,
		Password:  config.Password,
	})
	if err != nil {
		slog.Debug("etcd", "new client with password", err.Error())
		return nil, backErr.Wrap(err)
	}
	return cli, backErr.Submit()
}

// 解析bytes value值转成string
func (c *client) GetServerEntries(prefixKey string) ([]string, error) {
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
	if c.cli == nil {
		return backErr.Wrap(xerr.EtcdErrClientNotInitialized)
	}
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
	putResonse, err := c.kv.Put(
		c.ctx,
		s.Key,
		s.Value,
		clientv3.WithLease(c.leaseID),
	)
	// fmt.Print(putResonse)
	slog.Info("etcd", "action", "register", "put responce", putResonse)
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
