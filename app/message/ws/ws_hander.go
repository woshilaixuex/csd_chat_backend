package ws

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketHandler(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrader fail: %v", err)
		return
	}

	ctx := context.Background()
	client := &WSClient{
		WSID:     userID,
		Conn:     conn,
		Ctx:      ctx,
		Send:     make(chan []byte, 20),
		LastSeen: time.Now(),
	}

	DefaultClientManager.register <- client
	go client.readPump()
	go client.writePump()
}

func (c *WSClient) readPump() {
	defer func() {
		DefaultClientManager.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(maxLinkTime))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(maxLinkTime))
		c.LastSeen = time.Now()
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("连接异常: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

	}
}
func (c *WSClient) writePump() {
	ticker := time.NewTicker(maxLinkTime)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(maxLinkTime))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			sendmessage, err := json.Marshal(message)
			if err != nil {
				return
			}
			w.Write(sendmessage)

			// 批量发送积压消息
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(Line)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(maxLinkTime))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
