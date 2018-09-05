package pnet

import (
	"context"
	"net"
)

// ClientInfo 客户端信息结构体
type ClientInfo struct {
	Addr     string
	ClientID string
	Conn     net.Conn
	RW       *ReadWriter
}

// NewClientInfo 新建一个客户端信息结构
func NewClientInfo(clientID string, conn net.Conn) *ClientInfo {
	rw := newReaderWriterFromConn(clientID, conn)
	clientInfo := &ClientInfo{
		Addr:     conn.RemoteAddr().String(),
		ClientID: clientID,
		Conn:     conn,
		RW:       rw,
	}
	return clientInfo
}

// Write 向readwrite写入数据
func (ci *ClientInfo) Write(data []byte) error {
	return ci.RW.WritePack(data)
}

// ReadToMessageChan 监听消息
func (ci *ClientInfo) ReadToMessageChan(msgChan chan *Message) (context.Context, context.CancelFunc) {
	return ci.RW.ReadToMessageChan(msgChan)
}
