package main

import (
	"crypto/tls"
	"net"
)

func main() {
	host := "qiaoting"
	port := "10000"
	addr := net.JoinHostPort(host, port)

	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	_, err := tls.Dial("tcp", addr, config)
	if err != nil {
		panic(err)
	}
}
