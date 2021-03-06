package pnet

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"net"
	"time"

	"git.pleamon.com/p/plog"
)

// Client is Client struct.
type Client struct {
	Addr        string
	GetClientID func(net.Conn) string
	conn        net.Conn
	rw          *ReadWriter
	IsTLS       bool
	CACert      []byte
	PubKey      []byte
	PriKey      []byte
}

// NewClient create a new client.
func NewClient(addr string) (client *Client) {
	client = &Client{
		Addr:  addr,
		IsTLS: false,
	}
	return
}

// NewTLSClient create a new client.
func NewTLSClient(addr string, caCert, pubKey, priKey []byte) (client *Client) {
	client = &Client{
		Addr:   addr,
		IsTLS:  true,
		CACert: caCert,
		PubKey: pubKey,
		PriKey: priKey,
	}
	return
}

// Connect is Client connect to Server.
func (c *Client) Connect(timeout time.Duration) (err error) {
	dialer := &net.Dialer{
		Timeout: timeout,
	}
	var conn net.Conn
	if c.IsTLS {
		pool := x509.NewCertPool()
		ok := pool.AppendCertsFromPEM(c.CACert)
		if !ok {
			panic(err)
		}
		cer, err := tls.X509KeyPair(c.PubKey, c.PriKey)
		if err != nil {
			panic(err)
		}
		config := &tls.Config{
			RootCAs:      pool,
			Certificates: []tls.Certificate{cer},
		}
		conn, err = tls.DialWithDialer(dialer, "tcp", c.Addr, config)
		if err != nil {
			plog.Fatal(err)
			return err
		}
	} else {
		var err error
		plog.Debug("start connect")
		conn, err = dialer.Dial("tcp", c.Addr)
		if err != nil {
			plog.Fatal(err)
			return err
		}
		plog.Debug("connect success")
	}
	c.conn = conn
	c.rw = newReaderWriterFromConn(c.GetClientID(c.conn), c.conn)
	return
}

// Disconnect 断开连接
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
		plog.Debug("读取完成, ", err)
		return nil, err
	case err != nil:
		plog.Debug("读取出错, ", err)
		return nil, err
	}

	return msg, nil
}

// ReadToMessageChan 监听消息
func (c *Client) ReadToMessageChan(msgChan chan *Message) (ctx context.Context, cancel context.CancelFunc) {
	return c.rw.ReadToMessageChan(msgChan)
}

// Close 关闭连接
func (c *Client) Close() {
	c.conn.Close()
}
