package pnet

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
)

type HandleFunc func([]byte) ([]byte, error)
type Encoding interface {
	Encode([]byte) ([]byte, error, bool)
	Decode([]byte) ([]byte, error, bool)
}

type DefaultEncoding struct {
}

func (de *DefaultEncoding) Encode(data []byte) ([]byte, error, bool) {
	return data, nil, false
}

func (de *DefaultEncoding) Decode(data []byte) ([]byte, error, bool) {
	return data, nil, false
}

type EndPoint struct {
	handles map[string]HandleFunc
	keys    []string
}

type Server struct {
	Host          string
	Port          string
	Initinize     func(*Server)
	FinishConn    func(string, string, error)
	NoMatchHandle func(string, string, []byte) bool
	EndPoint      *EndPoint
	Encoding      Encoding
}

func NewServer(host, port string, e Encoding) *Server {
	endPoint := &EndPoint{
		handles: make(map[string]HandleFunc),
	}
	server := &Server{
		Host:     host,
		Port:     port,
		EndPoint: endPoint,
		Encoding: e,
	}
	return server
}

func (s *Server) AddHandle(name string, f HandleFunc) {
	s.EndPoint.handles[name] = f
	s.EndPoint.keys = append(s.EndPoint.keys, name)
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
		s.handleConn(conn)
	}
	return nil
}

func (s *Server) handleConn(conn net.Conn) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	rw := bufio.NewReadWriter(reader, writer)
	host, port, _ := net.SplitHostPort(conn.LocalAddr().String())
	for {
		line, _, err := rw.ReadLine()
		switch {
		case err == io.EOF:
			log.Println("读取完成, ", err.Error())
			if s.FinishConn != nil {
				s.FinishConn(host, port, err)
			}
			return
		case err != nil:
			log.Println("读取出错, ", err.Error())
			if s.FinishConn != nil {
				s.FinishConn(host, port, err)
			}
			return
		}

		data, err, close := s.Encoding.Decode(line)
		if err != nil {
			log.Println("解码失败, ", err.Error())
			if close {
				conn.Close()
				return
			}
			continue
		}
		var handle HandleFunc
		var handleName []byte
		for _, k := range s.EndPoint.keys {
			if bytes.HasPrefix(data, []byte(k)) {
				handle = s.EndPoint.handles[k]
				handleName = []byte(k)
				break
			}
		}

		rawData := bytes.TrimPrefix(data, handleName)
		if handle == nil {
			if s.NoMatchHandle != nil {
				if s.NoMatchHandle(host, port, rawData) {
					conn.Close()
					return
				}
			}
			continue
		}

		response, err := handle(rawData)
		if err != nil {
			conn.Close()
			return
		}

		decoded, err, close := s.Encoding.Encode(response)
		if err != nil {
			log.Println("编码失败, ", err.Error())
			if close {
				conn.Close()
				return
			}
			continue
		}
		rw.Write(decoded)
		rw.Flush()
	}
}
