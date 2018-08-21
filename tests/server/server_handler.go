package main

import (
	"log"

	"git.pleamon.com/p/pnet"
)

type ServerHandler struct {
	ClientInfo *pnet.ClientInfo
}

func (sh *ServerHandler) Initinize(s *pnet.Server) {
	log.Printf("启动服务程序成功，监听端口 [%s]\n", s.Addr)
}

func (sh *ServerHandler) AcceptConnHandle(s *pnet.Server, ci *pnet.ClientInfo) ([]byte, error) {
	sh.ClientInfo = ci
	log.Println("有新客户端加入")
	return nil, nil
}

func (sh *ServerHandler) HeathHandle(s *pnet.Server) error {
	log.Println("HeathHandle")
	return nil
}

func (sh *ServerHandler) ReceiveMessageHandle(s *pnet.Server, msg *pnet.Message) ([]byte, error) {
	log.Println(msg.Data)

	return []byte("hello world"), nil
}

func (sh *ServerHandler) FinishConnHandle(error) {
	log.Printf("disconnect %s", sh.ClientInfo.Addr)
}

func (sh *ServerHandler) ErrorHandle(int, interface{}) {

}
