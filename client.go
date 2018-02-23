package pnet

import (
	"context"
	"io"
	"log"
	"net"
)

// Client is Client struct.
type Client struct {
	Host        string
	Port        string
	GetClientID func(net.Conn) string
	conn        net.Conn
	rw          *ReadWriter
	Coding      *Coding
}

// NewClient create a new client.
func NewClient(host, port string) (client *Client) {
	client = &Client{
		Host: host,
		Port: port,
	}
	return
}

// Connect is Client connect to Server.
func (c *Client) Connect() (err error) {
	c.conn, err = net.Dial("tcp", net.JoinHostPort(c.Host, c.Port))
	if c.GetClientID == nil {
		c.GetClientID = GetClientID
	}
	if err == nil {
		c.rw = NewReaderWriterFromConn(c.GetClientID(c.conn), c.conn, c.Coding)
	}
	return
}

// Send send message to server
func (c *Client) Send(taskID, messageID uint64, dataBytes []byte) error {
	return c.rw.WritePack(taskID, messageID, dataBytes)
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
