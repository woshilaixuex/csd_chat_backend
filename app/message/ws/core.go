package ws

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 链接和管理核心
 * @Date: 2025-03-02 21:50
 */
var (
	maxClient      int   = 300              // 最大连接数
	maxLinkTimeNum       = 6                // 重试次数
	maxLinkTime          = time.Minute * 5  // 最大链接时长
	maxPingTime          = time.Second * 20 // 最大尝试时长
	maxMessageSize int64 = 512
	Line                 = []byte{'\n'}
	Space                = []byte{' '}
	Colon                = []byte{':'}
)

const (
	TEXT MType = 1 + iota
	IMAGE
)

type MType int

type Message struct {
	FromID  uint64 `json:"fromId"`
	SendID  uint64 `json:"sendId"`
	Content []byte `json:"content"`
	Type    MType  `json:"type"`
}

type WSClient struct {
	WSID     uint64
	Conn     *websocket.Conn
	Send     chan []byte
	LastSeen time.Time
}

func NewWSClinet(userID uint64, conn *websocket.Conn) *WSClient {
	return &WSClient{
		WSID:     userID,
		Conn:     conn,
		Send:     make(chan []byte, 20),
		LastSeen: time.Now(),
	}
}
func (ws *WSClient) Close() {
	if ws.Conn != nil {
		_ = ws.Conn.Close()
	}

	close(ws.Send)

}

type WSClientManager struct {
	clients    map[uint64]*WSClient
	chat       chan *Message // 广播通道
	register   chan *WSClient
	unregister chan *WSClient
	stop       chan struct{}
	c_num      int
	sync.RWMutex
}

func NewWSClientManager() *WSClientManager {
	return &WSClientManager{
		clients:    make(map[uint64]*WSClient),
		chat:       make(chan *Message),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
		stop:       make(chan struct{}),
		c_num:      0,
	}
}
