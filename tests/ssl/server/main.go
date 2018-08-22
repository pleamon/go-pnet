package main

import (
	"io/ioutil"
	"time"

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
	pubKey, err = ioutil.ReadFile("../ssl/server/server_cert.crt")
	if err != nil {
		panic(err)
	}
	priKey, err = ioutil.ReadFile("../ssl/server/server_key.pem")
	if err != nil {
		panic(err)
	}
}

func main() {
	server := pnet.NewTlsServer("127.0.0.1:10000", caCert, pubKey, priKey, &ServerHandler{}, time.Duration(10*time.Second))
	err := server.Listen()
	if err != nil {
		panic(err)
	}
}
