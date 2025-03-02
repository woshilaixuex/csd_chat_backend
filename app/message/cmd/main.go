package main

import (
	"github.com/gin-gonic/gin"
	"github.com/woshilaixuex/csd_chat_backend/app/message/router"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2025-02-27 14:18
 */

func main() {
	r := gin.Default()
	router.RouterRegister(r)
	r.Run()
}
