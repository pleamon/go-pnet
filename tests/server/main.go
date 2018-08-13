package main

import (
	"log"
	"net"

	"git.pleamon.com/p/pnet"
)

func Initinize(s *pnet.Server) {
	log.Printf("启动服务程序成功，监听端口 [%s]\n", s.Addr)
}

func AcceptConnHandle(conn net.Conn, rw *pnet.ReadWriter, clientId string) ([]byte, error) {
	return nil, nil
}

func MainHandle(server *pnet.Server, clientId *pnet.ClientInfo, msg *pnet.Message) ([]byte, error) {
	log.Println("client id: ", msg.ClientID)
	log.Println("length: ", msg.Length)
	log.Println("raw data:", msg.RawData, string(msg.RawData))
	log.Println("data: ", msg.Data, string(msg.Data))
	log.Println("this is server message")
	return nil, nil
}

func Encode(data []byte) []byte {
	log.Println("encode")
	return data
}

func Decode(data []byte) ([]byte, error) {
	log.Println("decode")
	return data, nil
}

func main() {
	// coding := &pnet.Coding{
	// Encode: Encode,
	// Decode: Decode,
	// }
	server := pnet.NewServer("127.0.0.1:10000", 0)
	server.Initinize = Initinize
	server.AsyncHandle = MainHandle
	// server.Coding = coding
	err := server.Listen()
	if err != nil {
		panic(err)
	}
}
