package main

import (
	"log"

	"git.pleamon.com/p/pnet"
)

func Initinize(s *pnet.Server) {
	log.Printf("启动服务程序成功，监听端口 [%s:%s]\n", s.Host, s.Port)
}

type Test struct {
	Property string
}

func (t *Test) Hello1(msg *pnet.Message) ([]byte, error) {
	log.Println(string(msg.Data))
	return nil, nil
}
func (t *Test) hello2() {

}

func MainHandle(msg *pnet.Message) (uint64, []byte, error) {
	log.Println("client id: ", msg.ClientId)
	log.Println("length: ", msg.Length)
	log.Println("task id: ", msg.TaskId)
	log.Println("raw data:", msg.RawData, string(msg.RawData))
	log.Println("data: ", msg.Data, string(msg.Data))
	return 100, []byte("this is server message"), nil
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
	coding := &pnet.Coding{
	// Encode: Encode,
	// Decode: Decode,
	}
	server := pnet.NewServer("127.0.0.1", "10000")
	server.Initinize = Initinize
	server.Handle = MainHandle
	server.Coding = coding
	err := server.Listen()
	if err != nil {
		panic(err)
	}
}
