package pnet

import (
	"crypto/tls"
	"log"
	"net"
	//"sync"
	"time"
)

type Server struct {
	Addr        string
	HeathTicker time.Duration
	//ClientPool       sync.Map
	ClientPool  map[string]*ClientInfo
	Cer         *tls.Certificate
	GetClientID func(net.Conn) string
	Handler     ServerHandler
	Coding      *Coding
}

func NewServer(addr string, handler ServerHandler, heathTickers ...time.Duration) *Server {
	server := &Server{
		Addr:       addr,
		Handler:    handler,
		ClientPool: make(map[string]*ClientInfo),
	}
	if len(heathTickers) > 0 {
		server.HeathTicker = heathTickers[0]
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
		log.Println("server listen error: ", err)
		return err
	}
	if s.Handler.Initinize != nil {
		s.Handler.Initinize(s)
	}
	if s.GetClientID == nil {
		s.GetClientID = GetClientID
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("accept: ", err)
		}
		go s.handleConn(conn)
	}
}

func (s *Server) createTicker(tick chan time.Time, done chan bool) {
	if s.HeathTicker == 0 {

	} else {
		ticker := time.NewTicker(s.HeathTicker)
		for {
			select {
			case <-ticker.C:
				tick <- time.Now()
			case <-done:
				ticker.Stop()
				break
			}
		}
	}
}

func (s *Server) handleConn(conn net.Conn) {
	clientID := s.GetClientID(conn)
	clientInfo := NewClientInfo(clientID, conn, s.Coding)
	if s.Coding != nil {
		clientInfo.RW.Coding = s.Coding
	}
	s.ClientPool[clientID] = clientInfo

	respData, err := s.Handler.AcceptConnHandle(s, clientInfo)
	if err != nil {
		// log.Println(err)
		conn.Close()
		return
	}
	clientInfo.Write(respData)

	msgChan := make(chan *Message, 100)

	ctx, cancel := clientInfo.ReadToMessageChan(msgChan)

	tick := make(chan time.Time)
	tickDone := make(chan bool, 1)
	go s.createTicker(tick, tickDone)
	for {
		select {
		case msg := <-msgChan:
			go func(ci *ClientInfo, msg *Message) {
				respData, err := s.Handler.ReceiveMessageHandle(s, msg)
				if err != nil {
					cancel()
					return
				}
				if len(respData) > 0 {
					ci.Write(respData)
				}
			}(clientInfo, msg)
		case <-ctx.Done():
			clientInfo.Lock()
			defer clientInfo.Unlock()
			s.Handler.FinishConnHandle(ctx.Err())
			delete(s.ClientPool, clientID)
			tickDone <- true
			conn.Close()
			return
		case <-tick:
			if err := s.Handler.HeathHandle(s); err != nil {
				cancel()
			}
		}
	}
}
