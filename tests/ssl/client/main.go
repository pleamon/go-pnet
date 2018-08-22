package main

import (
	"io/ioutil"
	"log"
	"net"
	"sync"

	"git.pleamon.com/p/pnet"
)

var (
	caCert []byte
	pubKey []byte
	priKey []byte
)

func init() {
	var err error
	caCert, err = ioutil.ReadFile("../ssl/ca/cacert.pem")
	if err != nil {
		panic(err)
	}
	pubKey, err = ioutil.ReadFile("../ssl/client/client_cert.crt")
	if err != nil {
		panic(err)
	}
	priKey, err = ioutil.ReadFile("../ssl/client/client_key.pem")
	if err != nil {
		panic(err)
	}
}

func main() {
	var wg sync.WaitGroup
	msgChan := make(chan *pnet.Message, 100)
	host := "qiaoting"
	port := "10000"
	addr := net.JoinHostPort(host, port)
	client := pnet.NewTlsClient(addr, caCert, pubKey, priKey)
	err := client.Connect()
	if err != nil {
		panic(err)
	}

	ctx, _ := client.ReadToMessageChan(msgChan)

	go func() {
		for {
			select {
			case msg := <-msgChan:
				log.Println("client id: ", msg.ClientID)
				log.Println("length: ", msg.Length)
				log.Println("raw data:", msg.RawData, string(msg.RawData))
				log.Println("data: ", msg.Data, string(msg.Data))
			case <-ctx.Done():
				log.Println("done")
				return
			}
			wg.Done()
		}
	}()

	err = client.Send([]byte("this is message 1"))
	if err != nil {
		panic(err)
	}
	wg.Add(1)

	err = client.Send([]byte("this is message 2"))
	if err != nil {
		panic(err)
	}
	wg.Add(1)

	err = client.Send([]byte("this is message 3"))
	if err != nil {
		panic(err)
	}
	wg.Add(1)

	wg.Wait()
	client.Close()
}
