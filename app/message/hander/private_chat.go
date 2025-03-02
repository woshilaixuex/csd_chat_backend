package hander

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/woshilaixuex/csd_chat_backend/app/message/ws"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 私聊
 * @Date: 2025-02-27 21:10
 */
type PrivateChatRequest struct {
	SendID  int64  `json:"send_id" binding:"required"`
	Content string `json:"content" binding:"required"`
	Type    int    `json:"type"`
}

func PrivateChatHander(ctx *gin.Context) {
	// 1. 从上下文中获取发送者ID（经过中间件验证后的用户ID）
	senderID := ctx.MustGet("userID").(uint64)

	// 2. 解析请求体
	var req PrivateChatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// 3. 构造消息对象（适配 ws 包的定义）
	msg := &ws.Message{
		FromID:  int64(senderID), // 发送者ID（需类型转换）
		SendID:  req.SendID,      // 接收者ID（来自请求体）
		Content: []byte(req.Content),
		Type:    ws.MType(req.Type), // 假设 ws 包有对应的 MType 定义
	}

	// 4. 发送消息
	ws.DefaultClientManager.SendChat(msg)

	// 5. 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{"status": "message sent"})
}
