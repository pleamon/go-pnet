package main

import (
	"log"
	"sync"

	"git.pleamon.com/p/pnet"
)

func main() {
	var wg sync.WaitGroup
	msgChan := make(chan *pnet.Message, 100)
	host := "127.0.0.1"
	port := "10000"
	client := pnet.NewClient(host, port)
	client.Connect()

	ctx, _ := client.ReadToMessageChan(msgChan)

	go func() {
		for {
			select {
			case msg := <-msgChan:
				log.Println("client id: ", msg.ClientID)
				log.Println("length: ", msg.Length)
				log.Println("task id: ", msg.TaskID)
				log.Println("msg id: ", msg.MessageID)
				log.Println("raw data:", msg.RawData, string(msg.RawData))
				log.Println("data: ", msg.Data, string(msg.Data))
			case <-ctx.Done():
				log.Println("done")
				return
			}
			wg.Done()
		}
	}()

	msgId, err := client.Send(1, []byte("this is message 1"))
	if err != nil {
		panic(err)
	}
	log.Println("send 1 msg id:", msgId)
	wg.Add(1)

	msgId, err = client.Send(2, []byte("this is message 2"))
	if err != nil {
		panic(err)
	}
	log.Println("send 2 msg id:", msgId)
	wg.Add(1)

	msgId, err = client.Send(3, []byte("this is message 3"))
	if err != nil {
		panic(err)
	}
	log.Println("send 3 msg id:", msgId)
	wg.Add(1)

	wg.Wait()
	client.Close()
}
