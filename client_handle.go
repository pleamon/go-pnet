package pnet

import (
	"net"
)

// ClientHandler 客户端回调函数
type ClientHandler struct {
	GetClientID func(net.Conn) string
	// 程序初始化时调用
	Initinize func(*Server)
	// 有新客户端连接时调用
	OnAccept func(*Server, *ClientInfo) []byte
	// 心跳检测
	OnHeath func(*Server, *ClientInfo) error
	// 接收到客户端数据时调用
	OnReceive func(*Server, *ClientInfo, *Message) []byte
	// 客户端连接断开时调用
	OnClose func(*Server, *ClientInfo, interface{})
	// 程序错误时调用
	OnError func(*Server, *ClientInfo, error)
}

// NewClientHandler 创建客户端回调函数
func NewClientHandler(config *ServerConfig) *ClientHandler {
	return &ClientHandler{
		GetClientID: config.GetClientID,
		Initinize:   config.Initinize,
		OnAccept:    config.OnAccept,
		OnHeath:     config.OnHeath,
		OnReceive:   config.OnReceive,
		OnClose:     config.OnClose,
		OnError:     config.OnError,
	}
}
