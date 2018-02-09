package main

import (
	"log"

	"git.pleamon.com/p/pnet"
)

func main() {
	host := "127.0.0.1"
	port := "10000"
	client := pnet.NewClient(host, port)
	client.Connect()
	err := client.Send([]byte("hello world"))
	if err != nil {
		panic(err)
	}
	data, err := client.Read()
	if err != nil {
		panic(err)
	}
	log.Println("length: ", data.Length)
	log.Println("task id: ", data.TaskId)
	log.Println("raw data:", data.RawData, string(data.RawData))
	log.Println("data: ", data.Data, string(data.Data))
	client.Close()
}
