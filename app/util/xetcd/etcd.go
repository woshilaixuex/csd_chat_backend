package xetcd

import (
	"log"
	"sync"
	"time"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: etcd核心模块
 * @Date: 2025-03-19 20:40
 */

// etcd:默认信息
const (
	defaultRenewalTime time.Duration = 10 * time.Second
)

// etcd:全局
var (
	etcdOnce sync.Once
)

// 实现对etcd检测变化的应略
type Instancer struct {
	client Client
	prefix string
	logger log.Logger
	quitc  chan struct{}
}

// Stop terminates the Instancer.
func (s *Instancer) Stop() {
	close(s.quitc)
}
