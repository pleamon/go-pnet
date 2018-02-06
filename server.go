package pnet

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
)

type Server struct {
	Host             string
	Port             string
	Initinize        func(*Server)
	AcceptConnHandle func(net.Conn)
	Handle           func([]byte, int64) ([]byte, error)
	FinishConnHandle func(net.Conn, error)
}

func NewServer(host, port string) *Server {
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
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	rw := bufio.NewReadWriter(reader, writer)
	for {
		dataSizeByte := make([]byte, 8)
		_, err := rw.Read(dataSizeByte)
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

		dataSizeBuffer := bytes.NewBuffer(dataSizeByte)
		var dataLength int64
		binary.Read(dataSizeBuffer, binary.BigEndian, &dataLength)

		dataByte := make([]byte, dataLength)
		_, err = rw.Read(dataByte)
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

		respData, err := s.Handle(dataByte, dataLength)
		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}
		if len(respData) == 0 {
			continue
		}
		respLen := uint64(len(respData))
		respPackLen := make([]byte, 8)
		binary.BigEndian.PutUint64(respPackLen, respLen)

		buffer := &bytes.Buffer{}
		buffer.Write(respPackLen)
		buffer.Write(respData)

		rw.Write(buffer.Bytes())
		rw.Flush()
	}
}
