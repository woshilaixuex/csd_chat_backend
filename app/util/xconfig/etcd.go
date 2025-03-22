package xconfig

import "time"

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: etcd配置
 * @Date: 2025-03-16 21:45
 */

var EtcdConfigName = "etcd"

type EtcdConfig struct {
	Name      string
	Endpoints []string
	Time      time.Duration
	UserName  string
	Password  string
	Method    string
}

func NewEtcdConfig() *EtcdConfig {
	return new(EtcdConfig)
}

func (c *EtcdConfig) Bind() error {

	return nil
}

func (c *EtcdConfig) GetConfigName() string {
	return EtcdConfigName
}
