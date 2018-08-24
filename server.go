package pnet

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"sync"

	//"sync"
	"time"

	"git.pleamon.com/p/plog"
)

type Server struct {
	Addr        string
	HeathTicker time.Duration
	ClientPool  *ClientPool
	lock        *sync.Mutex
	GetClientID func(net.Conn) string
	Handler     ServerHandler
	Coding      *Coding
	IsTLS       bool
	CACert      []byte
	PubKey      []byte
	PriKey      []byte
}

func NewServer(addr string, handler ServerHandler, heathTickers ...time.Duration) *Server {
	server := &Server{
		Addr:       addr,
		Handler:    handler,
		ClientPool: NewClientPool(),
		IsTLS:      false,
	}
	if len(heathTickers) > 0 {
		server.HeathTicker = heathTickers[0]
	}
	return server
}

func NewTlsServer(addr string, caCert, pubKey, priKey []byte, handler ServerHandler, heathTickers ...time.Duration) *Server {
	server := &Server{
		Addr:       addr,
		Handler:    handler,
		ClientPool: NewClientPool(),
		IsTLS:      true,
		CACert:     caCert,
		PubKey:     pubKey,
		PriKey:     priKey,
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

	s.Handler.OnServerInitinize(s)
	if s.GetClientID == nil {
		s.GetClientID = GetClientID
	}
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
	clientID := s.GetClientID(conn)
	clientInfo := NewClientInfo(clientID, conn, s.Coding)
	if s.Coding != nil {
		clientInfo.RW.Coding = s.Coding
	}

	s.ClientPool.Set(clientID, clientInfo)

	respData, err := s.Handler.OnAccept(s, clientInfo)
	if err != nil {
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
				respData, err := s.Handler.OnReceive(s, msg)
				if err != nil {
					cancel()
					return
				}
				if len(respData) > 0 {
					ci.Write(respData)
				}
			}(clientInfo, msg)
		case <-ctx.Done():
			s.Handler.OnClose(ctx.Err())
			s.ClientPool.Del(clientID)
			tickDone <- true
			conn.Close()
			return
		case <-tick:
			if err := s.Handler.OnHeath(s); err != nil {
				cancel()
			}
		}
	}
}
