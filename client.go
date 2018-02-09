package pnet

import (
	"io"
	"log"
	"net"
)

type Client struct {
	Host   string
	Port   string
	conn   net.Conn
	rw     *ReadWriter
	Coding *Coding
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
	if err == nil {
		c.rw = NewReaderWriterFromConn(c.conn, c.Coding)
	}
	return
}

func (c *Client) Send(dataBytes []byte) error {
	return c.rw.WritePack(dataBytes)
}

func (c *Client) Read() (*Message, error) {
	msg, err := c.rw.ReadPack()
	switch {
	case err == io.EOF:
		log.Println("读取完成, ", err.Error())
		return nil, err
	case err != nil:
		log.Println("读取出错, ", err.Error())
		return nil, err
	}

	return msg, nil
}

func (c *Client) Close() {
	c.conn.Close()
}
