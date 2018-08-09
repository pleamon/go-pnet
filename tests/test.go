package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"log"
)

func RsaEncrypt(origData []byte) ([]byte, error) {
	publicKey, err := ioutil.ReadFile("public.pem")
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

func RsaDecrypt(ciphertext []byte) ([]byte, error) {
	privateKey, err := ioutil.ReadFile("private.pem")
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

func main() {
	e, err := RsaEncrypt([]byte("hello world"))
	if err != nil {
		panic(err)
	}
	log.Println(e)

	d, err := RsaDecrypt(e)
	if err != nil {
		panic(err)
	}
	log.Println(d)
	log.Println(string(d))
}
