package pnet

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
)

type Client struct {
	Host string
	Port string
	conn net.Conn
}

func NewClient(host, port string) (client *Client) {
	client = &Client{
		Host: host,
		Port: port,
	}
	return
}

func (c *Client) Connect() (err error) {
	c.conn, err = net.Dial("tcp", net.JoinHostPort(c.Host, c.Port))
	return
}

func (c *Client) Send(dataBytes []byte) error {
	dataLen := uint64(len(dataBytes))
	dataLenBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(dataLenBytes, dataLen)
	buffer := &bytes.Buffer{}
	buffer.Write(dataLenBytes)
	buffer.Write(dataBytes)
	_, err := c.conn.Write(buffer.Bytes())
	return err
}

func (c *Client) Read() ([]byte, error) {
	dataSizeByte := make([]byte, 8)
	_, err := c.conn.Read(dataSizeByte)
	switch {
	case err == io.EOF:
		log.Println("读取完成, ", err.Error())
		return nil, err
	case err != nil:
		log.Println("读取出错, ", err.Error())
		return nil, err
	}

	dataSizeBuffer := bytes.NewBuffer(dataSizeByte)
	var dataLength int64
	binary.Read(dataSizeBuffer, binary.BigEndian, &dataLength)

	dataByte := make([]byte, dataLength)
	_, err = c.conn.Read(dataByte)
	switch {
	case err == io.EOF:
		log.Println("读取完成, ", err.Error())
		return nil, err
	case err != nil:
		log.Println("读取出错, ", err.Error())
		return nil, err
	}
	return dataByte, nil
}

func (c *Client) Close() {
	c.conn.Close()
}
