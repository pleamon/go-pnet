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

func (t *Test) Hello1(data []byte) ([]byte, error) {
	log.Println(string(data))
	return nil, nil
}
func (t *Test) hello2() {

}

func MainHandle(data []byte, length int64) ([]byte, error) {
	log.Println(data)
	log.Println(string(data))
	return []byte("this is server message"), nil
}

func main() {
	server := pnet.NewServer("127.0.0.1", "10000")
	server.Initinize = Initinize
	server.Handle = MainHandle
	err := server.Listen()
	if err != nil {
		panic(err)
	}
}
