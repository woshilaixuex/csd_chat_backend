package xetcd

import (
	"log/slog"
	"sync"
	"time"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 服务端api 客户端基本功能:租约注册；set注册服务配置信息；租约续约功能；对注册配置信息的更新；
 * @Date: 2025-03-19 20:41
 */
const minHeartBeatTime = 500 * time.Millisecond

// Registrar registers service instance liveness information to etcd.
type Registrar struct {
	client  Client
	quitmtx sync.Mutex
	quit    chan struct{}
}

// Service holds the instance identifying data you want to publish to etcd. Key
// must be unique, and value is the string returned to subscribers, typically
// called the "instance" string in other parts of package sd.
type Service struct {
	Key   string // unique key, e.g. "/service/foobar/1.2.3.4:8080"
	Value string // returned to subscribers, e.g. "http://1.2.3.4:8080"
	TTL   *TTLOption
}

// TTLOption allow setting a key with a TTL. This option will be used by a loop
// goroutine which regularly refreshes the lease of the key.
type TTLOption struct {
	heartbeat time.Duration // e.g. time.Second * 3
	ttl       time.Duration // e.g. time.Second * 10
}

var defaultTTLoption = &TTLOption{
	heartbeat: time.Second * 3,
	ttl:       time.Second * 10,
}

// NewTTLOption returns a TTLOption that contains proper TTL settings. Heartbeat
// is used to refresh the lease of the key periodically; its value should be at
// least 500ms. TTL defines the lease of the key; its value should be
// significantly greater than heartbeat.
//
// Good default values might be 3s heartbeat, 10s TTL.
func NewTTLOption(heartbeat, ttl time.Duration) *TTLOption {
	if heartbeat <= minHeartBeatTime {
		heartbeat = minHeartBeatTime
	}
	if ttl <= heartbeat {
		ttl = 3 * heartbeat
	}
	return &TTLOption{
		heartbeat: heartbeat,
		ttl:       ttl,
	}
}

func NewRegistrar(client Client) *Registrar {
	return &Registrar{
		client: client,
	}
}

func (r *Registrar) Register(service Service) {
	if err := r.client.Register(service); err != nil {
		return
	}
	if service.TTL != nil {
		slog.Info("etcd", "action", "register", "lease", r.client.LeaseID())
	} else {
		slog.Info("etcd", "action", "register")
	}
}
func (r *Registrar) Deregister(service Service) {
	if err := r.client.Deregister(service); err != nil {
		slog.Error("etcd", "action", "deregister", "err", err.Error())
	} else {
		slog.Info("etcd", "action", "deregister")
	}

	r.quitmtx.Lock()
	defer r.quitmtx.Unlock()
	if r.quit != nil {
		close(r.quit)
		r.quit = nil
	}
}
