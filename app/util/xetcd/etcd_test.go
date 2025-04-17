package xetcd_test

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/woshilaixuex/csd_chat_backend/app/util/xetcd"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: etcd测试
 * @Date: 2025-04-14 23:37
 */
func TestClinet(t *testing.T) {
	client, err := xetcd.NewClinet(context.Background())
	if err != nil {
		slog.Error("etcd client", "err", err.Error())
		t.Fatal()
	}
	slog.Debug("etcd client", "leaseid", client.LeaseID())
}

func TestRegister(t *testing.T) {
	client, err := xetcd.NewClinet(context.Background())
	if err != nil {
		slog.Error("etcd client", "err", err.Error())
		t.Fatal()
	}
	server := xetcd.Service{
		Key:   "/sever",
		Value: "0.0.0.0",
	}
	fmt.Print()
	client.Register(server)
	time.Sleep(5 * time.Second)
}
