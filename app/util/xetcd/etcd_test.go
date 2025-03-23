package xetcd_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/woshilaixuex/csd_chat_backend/app/util/xetcd"
)

func TestClinet(t *testing.T) {
	client, err := xetcd.NewClinet(context.Background(), xetcd.ClientOptions{})

	if err != nil {
		slog.Error("etcd client", "err", err.Error())
		t.Fatal()
	}
	slog.Debug("etcd client", "leaseid", client.LeaseID())
}
