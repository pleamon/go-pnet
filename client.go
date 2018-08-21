package pnet

import (
	"context"
	"io"
	"log"
	"net"
)

// Client is Client struct.
type Client struct {
	Addr        string
	GetClientID func(net.Conn) string
	conn        net.Conn
	rw          *ReadWriter
	Coding      *Coding
}

// NewClient create a new client.
func NewClient(addr string) (client *Client) {
	client = &Client{
		Addr: addr,
	}
	return
}

// Connect is Client connect to Server.
func (c *Client) Connect() (err error) {
	conn, err := net.Dial("tcp", c.Addr)
	if err != nil {
		return err
	}
	c.conn = conn
	if c.GetClientID == nil {
		c.GetClientID = GetClientID
	}
	c.rw = NewReaderWriterFromConn(c.GetClientID(c.conn), c.conn, c.Coding)
	return
}

func (c *Client) Disconnect() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Send send message to server
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

func (c *Client) ReadToMessageChan(msgChan chan *Message) (ctx context.Context, cancel context.CancelFunc) {
	return c.rw.ReadToMessageChan(msgChan)
}

func (c *Client) Close() {
	c.conn.Close()
}
