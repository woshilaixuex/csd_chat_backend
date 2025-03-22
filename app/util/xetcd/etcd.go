package xetcd

import (
	"sync"
	"time"

	"google.golang.org/grpc"
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

type Instancer struct {
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

	Username string
	Password string
}
