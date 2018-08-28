package pnet

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"

	//"sync"
	"time"

	"git.pleamon.com/p/plog"
)

type Server struct {
	Addr          string
	HeathTicker   time.Duration
	ClientPool    *ClientPool
	initinize     func(*Server)
	Coding        *Coding
	IsTLS         bool
	CACert        []byte
	PubKey        []byte
	PriKey        []byte
	ClientHandler *ClientHandler
}

func NewServer(addr string, config *ServerConfig) *Server {
	server := &Server{
		Addr:          addr,
		ClientPool:    NewClientPool(),
		ClientHandler: NewClientHandler(config),
		HeathTicker:   config.HeathTicker,
	}
	if config.CACert != nil {
		server.IsTLS = true
		server.CACert = config.CACert
		server.PubKey = config.PubKey
		server.PriKey = config.PriKey
	}
	return server
}

func (s *Server) Listen() error {

	var ln net.Listener
	if s.IsTLS {
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(s.CACert)

		cer, err := tls.X509KeyPair(s.PubKey, s.PriKey)
		if err != nil {
			plog.Fatal(err.Error())
		}
		config := &tls.Config{
			Certificates: []tls.Certificate{cer},
			ClientCAs:    pool,
			ClientAuth:   tls.RequireAndVerifyClientCert,
		}
		ln, err = tls.Listen("tcp", s.Addr, config)
		if err != nil {
			return err
		}
	} else {
		var err error
		ln, err = net.Listen("tcp", s.Addr)
		if err != nil {
			log.Println("server listen error: ", err)
			return err
		}
	}
	s.ClientHandler.Initinize(s)

	for {
		conn, err := ln.Accept()
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
	clientID := s.ClientHandler.GetClientID(conn)
	clientInfo := NewClientInfo(clientID, conn, s.Coding)
	if s.Coding != nil {
		clientInfo.RW.Coding = s.Coding
	}

	s.ClientPool.Set(clientID, clientInfo)

	if data := s.ClientHandler.OnAccept(s, clientInfo); data != nil {
		clientInfo.Write(data)
	}

	msgChan := make(chan *Message, 100)

	ctx, cancel := clientInfo.ReadToMessageChan(msgChan)

	tick := make(chan time.Time)
	tickDone := make(chan bool, 1)
	go s.createTicker(tick, tickDone)
	for {
		select {
		case msg := <-msgChan:
			go func(ci *ClientInfo, msg *Message) {
				if data := s.ClientHandler.OnReceive(s, clientInfo, msg); data != nil {
					ci.Write(data)
				}
			}(clientInfo, msg)
		case <-ctx.Done():
			s.ClientHandler.OnClose(s, clientInfo, ctx.Err())
			s.ClientPool.Del(clientID)
			tickDone <- true
			conn.Close()
			return
		case <-tick:
			if err := s.ClientHandler.OnHeath(s, clientInfo); err != nil {
				cancel()
			}
		}
	}
}
