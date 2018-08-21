package main

import (
	"log"
	"net"
	"sync"

	"git.pleamon.com/p/pnet"
)

func main() {
	var wg sync.WaitGroup
	msgChan := make(chan *pnet.Message, 100)
	host := "127.0.0.1"
	port := "10000"
	addr := net.JoinHostPort(host, port)
	client := pnet.NewClient(addr)
	client.Connect()

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

	err := client.Send([]byte("this is message 1"))
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
