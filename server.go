package pnet

import (
	"crypto/tls"
	"io"
	"log"
	"net"
)

type Server struct {
	Host             string
	Port             string
	Cer              *tls.Certificate
	Initinize        func(*Server)
	AcceptConnHandle func(net.Conn)
	Handle           func(*Message) ([]byte, error)
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
		Host:   host,
		Port:   port,
		PubKey: pubKey,
		PriKey: priKey,
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
	rw := NewReaderWriterFromConn(conn, s.Coding)
	for {
		msg, err := rw.ReadPack()
		switch {
		case err == io.EOF:
			log.Println("读取完成, ", err.Error())
			if s.FinishConnHandle != nil {
				s.FinishConnHandle(conn, err)
			}
			return
		case err != nil:
			log.Println("读取出错, ", err.Error())
			if s.FinishConnHandle != nil {
				s.FinishConnHandle(conn, err)
			}
			return
		}

		if s.Handle != nil {
			respData, err := s.Handle(msg)
			if err != nil {
				log.Println(err)
				conn.Close()
				return
			}
			rw.WritePack(respData)
		}
	}
}
