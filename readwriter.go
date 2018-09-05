package pnet

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
	"net"

	"git.pleamon.com/p/plog"
)

const (
	// CHECKSUMFLAG 校验数据
	CHECKSUMFLAG = 1024
)

// ReadWriter 扩展RW结构体
type ReadWriter struct {
	ClientID  string
	MessageID uint32
	*bufio.Writer
	*bufio.Reader
	LenPlace      int
	CheckSumPlace int
}

// Message 回调消息体
type Message struct {
	ClientID string
	Length   uint32
	Data     []byte
	err      error
}

// newReaderWriterFromConn 从conn创建readwrite
func newReaderWriterFromConn(clientID string, conn net.Conn) *ReadWriter {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	rw := &ReadWriter{
		ClientID:      clientID,
		Reader:        reader,
		Writer:        writer,
		LenPlace:      4,
		CheckSumPlace: 4,
		MessageID:     0,
	}
	return rw
}

// ResetMessageID 重置消息ID
func (rw *ReadWriter) ResetMessageID() {
	rw.MessageID = 0
}

// ReadPack 读取一个包
func (rw *ReadWriter) ReadPack() (*Message, error) {
	dataLength, err := rw.readPackLen()
	if err != nil {
		plog.Error(err)
		return nil, err
	}
	checkSum, err := rw.readCheckSum()
	if err != nil {
		plog.Error(err)
	}
	if checkSum == 0 {
		plog.DebugF("checksum is %d, discard conn buffer data", checkSum)
		io.Copy(ioutil.Discard, rw)
		return nil, nil
	}
	if (dataLength+checkSum)%CHECKSUMFLAG != 0 {
		plog.DebugF("(checksum[%d] + datalength[%d]) % 1024 == 0, discard conn buffer data", checkSum, dataLength)
		io.Copy(ioutil.Discard, rw)
		return nil, nil
	}
	data, err := rw.readPackData(dataLength)
	if err != nil {
		plog.Error(err)
		return nil, ErrorCheckSumFailed
	}
	msg := &Message{
		ClientID: rw.ClientID,
		Length:   dataLength,
		Data:     data,
	}
	return msg, nil
}

// readPackLen 读取数据长度
func (rw *ReadWriter) readPackLen() (uint32, error) {
	dataSizeByte := make([]byte, rw.LenPlace)
	_, err := rw.Read(dataSizeByte)
	if err != nil {
		plog.Debug(err)
		return 0, err
	}
	dataSizeBuffer := bytes.NewBuffer(dataSizeByte)
	var dataLength uint32
	binary.Read(dataSizeBuffer, binary.BigEndian, &dataLength)

	return dataLength, nil
}

// readCheckSum 读取校验数据
func (rw *ReadWriter) readCheckSum() (uint32, error) {
	dataSizeByte := make([]byte, rw.CheckSumPlace)
	_, err := rw.Read(dataSizeByte)
	if err != nil {
		plog.Debug(err)
		return 0, err
	}
	dataSizeBuffer := bytes.NewBuffer(dataSizeByte)
	var dataCheckSum uint32
	binary.Read(dataSizeBuffer, binary.BigEndian, &dataCheckSum)

	return dataCheckSum, nil
}

// readPackData 读取数据区
func (rw *ReadWriter) readPackData(length uint32) ([]byte, error) {
	dataByte := make([]byte, length)
	_, err := rw.Read(dataByte)
	if err != nil {
		plog.Error(err)
		return nil, err
	}
	return dataByte, nil
}

// WritePack 向writer写入一个包
func (rw *ReadWriter) WritePack(dataByte []byte) error {
	length := len(dataByte)
	if length == 0 {
		return errors.New("not data")
	}
	dataLength := uint32(length)
	mold := length % CHECKSUMFLAG
	checkSum := uint32(CHECKSUMFLAG - mold)
	if checkSum == 0 {
		checkSum = CHECKSUMFLAG
	}
	encodeData := dataByte
	respPackLen := make([]byte, rw.LenPlace)
	binary.BigEndian.PutUint32(respPackLen, dataLength)

	respCheckSum := make([]byte, rw.CheckSumPlace)
	binary.BigEndian.PutUint32(respCheckSum, checkSum)

	buffer := &bytes.Buffer{}

	_, err := buffer.Write(respPackLen)
	if err != nil {
		plog.Error(err)
		return err
	}

	_, err = buffer.Write(respCheckSum)
	if err != nil {
		plog.Error(err)
		return err
	}

	_, err = buffer.Write(encodeData)
	if err != nil {
		plog.Error(err)
		return err
	}

	_, err = rw.Write(buffer.Bytes())
	if err != nil {
		plog.Error(err)
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
