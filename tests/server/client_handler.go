package main

import (
	"log"

	"git.pleamon.com/p/pnet"
)

type ClientHandler struct {
	Server     *pnet.Server
	ClientInfo *pnet.ClientInfo
}

func (sh *ClientHandler) OnAccept(server *pnet.Server, ci *pnet.ClientInfo) ([]byte, error) {
	sh.Server = server
	sh.ClientInfo = ci
	log.Println("有新客户端加入", ci.Addr)
	return nil, nil
}

func (sh *ClientHandler) OnHeath() error {
	log.Println("HeathHandle")
	return nil
}

func (sh *ClientHandler) OnReceive(msg *pnet.Message) ([]byte, error) {
	log.Println(msg.Data)

	return []byte("hello world"), nil
}

func (sh *ClientHandler) OnClose(error) {
	log.Printf("disconnect %s", sh.ClientInfo.Addr)
}

func (sh *ClientHandler) OnError(int, interface{}) {

}
