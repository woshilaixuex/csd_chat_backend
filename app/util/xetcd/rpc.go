package xetcd

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: rpc服务信息
 * @Date: 2025-03-16 21:38
 */

type ServiceMeta struct {
	Addr    string `json:"addr"`
	Port    int    `json:"port"`
	Type    string `json:"type"`
	Version string `json:"version"`
	Weight  int    `json:"weight"`
}
