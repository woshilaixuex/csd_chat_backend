package xetcd

import (
	"log"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 客户端api 客户端基本功能:get获取去注册信息；watch监听并更新服务端信息；
 * @Date: 2025-03-19 20:41
 */

func NewEtcdClient() *EtcdClinet {
	cli, err := newEtcdCon()
	if err != nil {
		log.Fatalln(err.Error())
		return nil
	}

}
