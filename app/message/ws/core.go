package ws

import (
	"context"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2025-03-02 21:50
 */
var (
	maxClient      int   = 300             // 最大连接数
	maxLinkTimeNum       = 6               // 重试次数
	maxLinkTime          = time.Minute * 5 // 最大重试时长
	maxPingTime          = time.Second * 20
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
	FromID  int64
	SendID  int64
	Content []byte
	Type    MType
}

type WSClient struct {
	WSID     int64
	Conn     *websocket.Conn
	Ctx      context.Context
	Send     chan []byte
	LastSeen time.Time
}

func (ws *WSClient) Close() {

}

type WSClientManager struct {
	clients    map[int64]*WSClient
	chat       chan *Message // 广播通道
	register   chan *WSClient
	unregister chan *WSClient
	c_num      int
	sync.RWMutex
}

func NewWSClientManager() *WSClientManager {
	return &WSClientManager{
		clients:    make(map[int64]*WSClient),
		chat:       make(chan *Message),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
		c_num:      0,
	}
}
