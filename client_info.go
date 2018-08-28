package pnet

import (
	"context"
	"net"
)

type ClientInfo struct {
	Addr     string
	ClientID string
	Conn     net.Conn
	RW       *ReadWriter
}

func NewClientInfo(clientID string, conn net.Conn) *ClientInfo {
	rw := NewReaderWriterFromConn(clientID, conn)
	clientInfo := &ClientInfo{
		Addr:     conn.RemoteAddr().String(),
		ClientID: clientID,
		Conn:     conn,
		RW:       rw,
	}
	return clientInfo
}

func GetClientID(conn net.Conn) string {
	return conn.RemoteAddr().String()
}

func (ci *ClientInfo) Write(data []byte) error {
	return ci.RW.WritePack(data)
}

func (ci *ClientInfo) ReadToMessageChan(msgChan chan *Message) (context.Context, context.CancelFunc) {
	return ci.RW.ReadToMessageChan(msgChan)
}
