package pnet

import (
	"net"
	"time"
)

type ServerConfig struct {
	CACert      []byte
	PubKey      []byte
	PriKey      []byte
	HeathTicker time.Duration
	GetClientID func(net.Conn) string
	Initinize   func(*Server)                               // 程序初始化时调用
	OnAccept    func(*Server, *ClientInfo) []byte           // 有新客户端连接时调用
	OnHeath     func(*Server, *ClientInfo) error            // 心跳检测
	OnReceive   func(*Server, *ClientInfo, *Message) []byte // 接收到客户端数据时调用
	OnError     func(*Server, *ClientInfo, error)           // 程序错误时调用
	OnClose     func(*Server, *ClientInfo, interface{})     // 客户端连接断开时调用
}
