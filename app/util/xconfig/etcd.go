package xconfig

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: etcd配置
 * @Date: 2025-03-16 21:45
 */

var EtcdConfigName = "etcd"

const (
	WITHUSERSTL = 1 << iota
	WITHUSERPASSWORD
)

type EtcdConfig struct {
	Name      string
	Endpoints []string
	Time      time.Duration
	UserName  string
	Password  string
	Method    int
}

func NewEtcdConfig() *EtcdConfig {
	return new(EtcdConfig)
}

func (c *EtcdConfig) Bind() error {
	c.Name = viper.GetString("etcd.name")
	endpoints := viper.GetString("etcd.endpoints")
	c.Endpoints = strings.Split(endpoints, ";")
	c.UserName = viper.GetString("etcd.username")
	c.Password = viper.GetString("etcd.password")
	c.Method = viper.GetInt("etcd.method")
	return nil
}

func (c *EtcdConfig) GetConfigName() string {
	return EtcdConfigName
}
