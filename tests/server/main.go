package main

import (
	"time"

	"git.pleamon.com/p/pnet"
)

func main() {
	// coding := &pnet.Coding{
	// Encode: Encode,
	// Decode: Decode,
	// }
	server := pnet.NewServer("127.0.0.1:10000", &ServerHandler{}, time.Duration(10*time.Second))
	err := server.Listen()
	if err != nil {
		panic(err)
	}
}
