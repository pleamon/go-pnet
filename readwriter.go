package pnet

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"

	"git.pleamon.com/p/plog"
)

const (
	CHECKSUMFLAG = 1024
)

type ReadWriter struct {
	ClientId  string
	MessageID uint32
	*bufio.Writer
	*bufio.Reader
	LenPlace      int
	CheckSumPlace int
}

type Message struct {
	ClientID string
	Length   uint32
	Data     []byte
	err      error
}

func NewReaderWriterFromConn(clientId string, conn net.Conn) *ReadWriter {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	rw := &ReadWriter{
		ClientId:      clientId,
		Reader:        reader,
		Writer:        writer,
		LenPlace:      4,
		CheckSumPlace: 4,
		MessageID:     0,
	}
	return rw
}

func (rw *ReadWriter) ResetMessageId() {
	rw.MessageID = 0
}

func (rw *ReadWriter) ReadPack() (*Message, error) {
	dataLength, err := rw.ReadPackLen()
	if err != nil {
		plog.Error(err)
		return nil, err
	}
	checkSum, err := rw.ReadCheckSum()
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
	data, err := rw.ReadPackData(dataLength)
	if err != nil {
		plog.Error(err)
		return nil, ErrorCheckSumFailed
	}
	msg := &Message{
		ClientID: rw.ClientId,
		Length:   dataLength,
		Data:     data,
	}
	return msg, nil
}

func (rw *ReadWriter) ReadPackLen() (uint32, error) {
	dataSizeByte := make([]byte, rw.LenPlace)
	_, err := rw.Read(dataSizeByte)
	if err != nil {
		plog.Debug(err)
		return 0, err
	}
	dataSizeBuffer := bytes.NewBuffer(dataSizeByte)
	var dataLength uint32
	binary.Read(dataSizeBuffer, binary.BigEndian, &dataLength)
	log.Println("datalength", dataLength)

	return dataLength, nil
}

func (rw *ReadWriter) ReadCheckSum() (uint32, error) {
	dataSizeByte := make([]byte, rw.CheckSumPlace)
	_, err := rw.Read(dataSizeByte)
	if err != nil {
		plog.Debug(err)
		return 0, err
	}
	dataSizeBuffer := bytes.NewBuffer(dataSizeByte)
	var dataCheckSum uint32
	binary.Read(dataSizeBuffer, binary.BigEndian, &dataCheckSum)
	log.Println("dataCheckSum", dataCheckSum)

	return dataCheckSum, nil
}

func (rw *ReadWriter) ReadPackData(length uint32) ([]byte, error) {
	dataByte := make([]byte, length)
	_, err := rw.Read(dataByte)
	if err != nil {
		plog.Error(err)
		return nil, err
	}
	return dataByte, nil
}

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
