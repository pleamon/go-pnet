package pnet

import (
	"crypto/tls"
	"log"
	"net"
)

type Server struct {
	Host             string
	Port             string
	Cer              *tls.Certificate
	GetClientID      func(net.Conn) string
	Initinize        func(*Server)
	AcceptConnHandle func(net.Conn)
	AsyncHandle      func(*Message) (uint64, []byte, error)
	SyncHandle       func(*Message) (uint64, []byte, error)
	FinishConnHandle func(net.Conn, error)
	Coding           *Coding
}

func NewServer(host, port string) *Server {
	server := &Server{
		Host: host,
		Port: port,
	}
	return server
}

func NewTlsServer(host, port, pubKey, priKey string) *Server {
	server := &Server{
		Host: host,
		Port: port,
	}
	return server
}

func (s *Server) Listen() error {
	l, err := net.Listen("tcp", net.JoinHostPort(s.Host, s.Port))
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
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	if s.AcceptConnHandle != nil {
		s.AcceptConnHandle(conn)
	}
	rw := NewReaderWriterFromConn(s.GetClientID(conn), conn, s.Coding)

	msgChan := make(chan *Message, 100)

	ctx, _ := rw.ReadToMessageChan(msgChan)

	for {
		select {
		case msg := <-msgChan:
			log.Println("msg: ", msg)
			if s.AsyncHandle != nil {
				go func(rw *ReadWriter) {
					respTaskID, respData, err := s.AsyncHandle(msg)
					if err != nil {
						log.Println(err)
						conn.Close()
						return
					}
					rw.WritePack(respTaskID, respData)
				}(rw)
				log.Println("next")
			}
		case <-ctx.Done():
			log.Println("err", ctx.Err())
			if s.FinishConnHandle != nil {
				s.FinishConnHandle(conn, ctx.Err())
			}
			return
		}
	}
}
