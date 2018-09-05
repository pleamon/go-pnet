package pnet

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net"

	"git.pleamon.com/p/plog"
)

const (
	BLOCKSIZE = 4096
)

var (
	ErrorCheckSumFailed = errors.New("error check sum failed")
)

type ReadWriter struct {
	ClientId  string
	MessageID uint64
	*bufio.Writer
	*bufio.Reader
	LenPlace int
}

type Message struct {
	ClientID string
	Length   int64
	Data     []byte
	err      error
}

func NewReaderWriterFromConn(clientId string, conn net.Conn) *ReadWriter {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	rw := &ReadWriter{
		ClientId:  clientId,
		Reader:    reader,
		Writer:    writer,
		LenPlace:  8,
		MessageID: 0,
	}
	return rw
}

func (rw *ReadWriter) ResetMessageId() {
	rw.MessageID = 0
}

func (rw *ReadWriter) ReadPack() (*Message, error) {
	dataLength, err := rw.ReadPackLen()
	msg := &Message{
		ClientID: rw.ClientId,
	}
	if err != nil {
		plog.Debug(err)
		return nil, err
	}
	if dataLength%BLOCKSIZE != 0 {
		// 丢弃数据
		io.Copy(ioutil.Discard, rw)
		return nil, errors.New("check sum failed")
	}
	data, err := rw.ReadPackData(dataLength)
	if err != nil {
		plog.Debug(err)
		return nil, ErrorCheckSumFailed
	}
	msg.Length = dataLength
	msg.Data = data
	return msg, nil
}

func (rw *ReadWriter) ReadPackLen() (int64, error) {
	dataSizeByte := make([]byte, rw.LenPlace)
	_, err := rw.Read(dataSizeByte)
	if err != nil {
		plog.Debug(err)
		return 0, err
	}
	dataSizeBuffer := bytes.NewBuffer(dataSizeByte)
	var dataLength int64
	binary.Read(dataSizeBuffer, binary.BigEndian, &dataLength)

	return dataLength, nil
}

func (rw *ReadWriter) ReadPackData(length int64) ([]byte, error) {
	if length > math.MaxInt64 {
		plog.DebugF("read pack data out of range [%d], current system max length [%d]", length, math.MaxInt64)
		return nil, fmt.Errorf("read pack data out of range [%d], current system max length [%d]", length, math.MaxInt64)
	}
	buffer := &bytes.Buffer{}
	for {
		if length <= 0 {
			break
		}
		dataByte := make([]byte, BLOCKSIZE)
		n, err := rw.Read(dataByte)
		if err != nil {
			plog.Error(err)
			return nil, err
		}
		buffer.Write(dataByte[0:n])
		length = length - BLOCKSIZE
	}
	return buffer.Bytes(), nil
}

func (rw *ReadWriter) WritePack(dataByte []byte) error {
	if len(dataByte) == 0 {
		return errors.New("not data")
	}
	length := len(dataByte)
	_length := length + (BLOCKSIZE - length%BLOCKSIZE)
	dataLength := uint64(_length)
	encodeData := dataByte
	respPackLen := make([]byte, rw.LenPlace)
	binary.BigEndian.PutUint64(respPackLen, dataLength)

	buffer := &bytes.Buffer{}

	_, err := buffer.Write(respPackLen)
	if err != nil {
		plog.Debug(err)
		return err
	}

	_, err = buffer.Write(encodeData)
	if err != nil {
		plog.Debug(err)
		return err
	}

	_, err = rw.Write(buffer.Bytes())
	if err != nil {
		plog.Debug(err)
		return err
	}
	return rw.Flush()
}

// ReadToMessageChan 读取客户端数据发送到管道
func (rw *ReadWriter) ReadToMessageChan(msgChan chan *Message) (ctx context.Context, cancel context.CancelFunc) {
	ctx, cancel = context.WithCancel(context.Background())
	go func() {
		for {
			msg, err := rw.ReadPack()
			if err == ErrorCheckSumFailed {
				continue
			}
			if err != nil {
				cancel()
				plog.Debug(err)
				return
			}
			msgChan <- msg
		}
	}()
	return ctx, cancel
}
