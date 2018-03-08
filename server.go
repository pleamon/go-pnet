package pnet

import (
	"crypto/tls"
	"log"
	"net"
)

type ClientInfo struct {
	ClientID string
	Conn     *net.Conn
	RW       *ReadWriter
}

var (
	ClientPool map[string]*ClientInfo
)

type Server struct {
	Addr             string
	Cer              *tls.Certificate
	GetClientID      func(*net.Conn) string
	Initinize        func(*Server)
	AcceptConnHandle func(*net.Conn, *ReadWriter, string) ([]byte, error)
	AsyncHandle      func(*Message) ([]byte, error)
	SyncHandle       func(*Message) ([]byte, error)
	FinishConnHandle func(*net.Conn, error)
	Coding           *Coding
}

func init() {
	ClientPool = make(map[string]*ClientInfo)
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
	ClientPool[clientID] = &ClientInfo{
		ClientID: clientID,
		Conn:     conn,
		RW:       rw,
	}
	if s.AcceptConnHandle != nil {
		respData, err := s.AcceptConnHandle(conn, rw, s.GetClientID(conn))
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
					respData, err := s.AsyncHandle(msg)
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
				s.FinishConnHandle(conn, ctx.Err())
			}
			delete(ClientPool, clientID)
			return
		}
	}
}
