package pnet

import (
	"net"
	"time"
)

// ServerConfig 服务器配置参数
type ServerConfig struct {
	CACert      []byte                                      // ca证书
	PubKey      []byte                                      // public证书
	PriKey      []byte                                      // private证书
	HeathTicker time.Duration                               // 心跳间隔 客户端健康检查
	GetClientID func(net.Conn) string                       // 获取客户端id
	Initinize   func(*Server)                               // 程序初始化时调用
	OnAccept    func(*Server, *ClientInfo) []byte           // 有新客户端连接时调用
	OnHeath     func(*Server, *ClientInfo) error            // 心跳检测
	OnReceive   func(*Server, *ClientInfo, *Message) []byte // 接收到客户端数据时调用
	OnError     func(*Server, *ClientInfo, error)           // 程序错误时调用
	OnClose     func(*Server, *ClientInfo, interface{})     // 客户端连接断开时调用
}
