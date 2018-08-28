package main

import (
	"log"
	"net"
	"time"

	"git.pleamon.com/p/pnet"
)

func initinize(server *pnet.Server) {
	log.Printf("启动服务程序成功，监听端口 [%s]", server.Addr)
}

func getClientId(conn net.Conn) string {
	return conn.RemoteAddr().String()
}

func onAccept(server *pnet.Server, clientInfo *pnet.ClientInfo) []byte {
	log.Println("有新客户端连接", clientInfo.Addr)
	return nil
}

func onHeath(server *pnet.Server, clientInfo *pnet.ClientInfo) error {
	log.Println("heath", clientInfo.Addr)
	return nil
}

func onReceive(server *pnet.Server, clientInfo *pnet.ClientInfo, message *pnet.Message) []byte {
	log.Println("on receive:", clientInfo.Addr)
	log.Printf("message %+v", message)
	return []byte("hello world")
}

func onError(server *pnet.Server, clientInfo *pnet.ClientInfo, err error) {
	log.Println("on error", clientInfo.Addr, err.Error())
}

func onClose(server *pnet.Server, clientInfo *pnet.ClientInfo, closeInfo interface{}) {
	log.Println("on close", clientInfo.Addr, closeInfo)
}

func main() {
	serverConfig := &pnet.ServerConfig{
		HeathTicker: time.Second * 10,
		GetClientID: getClientId,
		Initinize:   initinize,
		OnAccept:    onAccept,
		OnHeath:     onHeath,
		OnReceive:   onReceive,
		OnError:     onError,
		OnClose:     onClose,
	}
	server := pnet.NewServer("127.0.0.1:10000", serverConfig)
	err := server.Listen()
	if err != nil {
		panic(err)
	}
}
