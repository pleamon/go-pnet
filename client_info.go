package pnet

import (
	"context"
	"net"
	"sync"
)

type ClientInfo struct {
	Addr     string
	ClientID string
	Conn     net.Conn
	RW       *ReadWriter
	lock     *sync.Mutex
}

func NewClientInfo(clientID string, conn net.Conn, coding *Coding) *ClientInfo {
	rw := NewReaderWriterFromConn(clientID, conn, coding)
	clientInfo := &ClientInfo{
		Addr:     conn.RemoteAddr().String(),
		ClientID: clientID,
		Conn:     conn,
		RW:       rw,
		lock:     &sync.Mutex{},
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

func (ci *ClientInfo) Lock() {
	ci.lock.Lock()
}

func (ci *ClientInfo) Unlock() {
	ci.lock.Unlock()
}
