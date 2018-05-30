package pnet

import (
	"crypto/tls"
	"log"
	"net"
	"sync"
)

type ClientInfo struct {
	ClientID string
	Conn     *net.Conn
	RW       *ReadWriter
}

var ()

type Server struct {
	Addr             string
	ClientPool       sync.Map
	Cer              *tls.Certificate
	GetClientID      func(*net.Conn) string
	Initinize        func(*Server)
	AcceptConnHandle func(*Server, *ClientInfo) ([]byte, error)
	AsyncHandle      func(*Server, *ClientInfo, *Message) ([]byte, error)
	SyncHandle       func(*Message) ([]byte, error)
	FinishConnHandle func(string, *net.Conn, error)
	Coding           *Coding
}

func init() {
}

func NewServer(addr string) *Server {
	server := &Server{
		Addr: addr,
	}
	return server
}

func NewTlsServer(addr, pubKey, priKey string) *Server {
	server := &Server{
		Addr: addr,
	}
	return server
}

func (s *Server) Listen() error {
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	if s.Initinize != nil {
		s.Initinize(s)
	}
	if s.GetClientID == nil {
		s.GetClientID = GetClientID
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
		}
		go s.handleConn(&conn)
	}
}

func (s *Server) handleConn(conn *net.Conn) {
	clientID := s.GetClientID(conn)
	rw := NewReaderWriterFromConn(s.GetClientID(conn), conn, s.Coding)
	clientInfo := &ClientInfo{
		ClientID: clientID,
		Conn:     conn,
		RW:       rw,
	}
	s.ClientPool.Store(clientID, clientInfo)
	if s.AcceptConnHandle != nil {
		respData, err := s.AcceptConnHandle(s, clientInfo)
		if err != nil {
			log.Println(err)
			(*conn).Close()
			return
		}
		rw.WritePack(respData)
	}

	msgChan := make(chan *Message, 100)

	ctx, _ := rw.ReadToMessageChan(msgChan)

	for {
		select {
		case msg := <-msgChan:
			log.Println("msg: ", msg)
			if s.AsyncHandle != nil {
				go func(rw *ReadWriter, msg *Message) {
					respData, err := s.AsyncHandle(s, clientInfo, msg)
					if err != nil {
						log.Println(err)
						(*conn).Close()
						return
					}
					rw.WritePack(respData)
				}(rw, msg)
			}
		case <-ctx.Done():
			if s.FinishConnHandle != nil {
				s.FinishConnHandle(clientID, conn, ctx.Err())
			}
			s.ClientPool.Delete(clientID)
			return
		}
	}
}
