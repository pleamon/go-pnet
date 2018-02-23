package pnet

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"net"
)

type Coding struct {
	Encode func([]byte) []byte
	Decode func([]byte) ([]byte, error)
}

type ReadWriter struct {
	ClientId  string
	MessageID uint64
	*bufio.Writer
	*bufio.Reader
	LenPlace       int
	TaskIdPlace    int
	MessageIDPlace int
	Coding         *Coding
}

type Message struct {
	ClientID  string
	Length    int64
	TaskID    uint64
	MessageID uint64
	RawData   []byte
	Data      []byte
	err       error
}

func NewReaderWriterFromConn(clientId string, conn net.Conn, coding *Coding) *ReadWriter {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	rw := &ReadWriter{
		ClientId:       clientId,
		Reader:         reader,
		Writer:         writer,
		LenPlace:       8,
		TaskIdPlace:    8,
		MessageIDPlace: 8,
		MessageID:      0,
		Coding:         coding,
	}
	return rw
}

func GetClientID(conn net.Conn) string {
	return conn.RemoteAddr().String()
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
		return nil, err
	}
	dataTaskID, err := rw.ReadPackTaskID()
	if err != nil {
		return nil, err
	}
	dataMessageID, err := rw.ReadMessageID()
	if err != nil {
		return nil, err
	}
	data, err := rw.ReadPackData(dataLength)
	if err != nil {
		return nil, err
	}
	msg.Length = dataLength
	msg.TaskID = dataTaskID
	msg.MessageID = dataMessageID
	msg.RawData = data
	if rw.Coding != nil && rw.Coding.Decode != nil {
		msg.Data, msg.err = rw.Coding.Decode(msg.RawData)
	} else {
		msg.Data = msg.RawData
	}
	return msg, nil
}

func (rw *ReadWriter) ReadPackLen() (int64, error) {
	dataSizeByte := make([]byte, rw.LenPlace)
	_, err := rw.Read(dataSizeByte)
	if err != nil {
		return 0, err
	}
	dataSizeBuffer := bytes.NewBuffer(dataSizeByte)
	var dataLength int64
	binary.Read(dataSizeBuffer, binary.BigEndian, &dataLength)

	return dataLength, nil
}

func (rw *ReadWriter) ReadPackTaskID() (uint64, error) {
	dataByte := make([]byte, rw.TaskIdPlace)
	_, err := rw.Read(dataByte)
	if err != nil {
		return 0, err
	}
	dataBuffer := bytes.NewBuffer(dataByte)
	var dataTaskID uint64
	binary.Read(dataBuffer, binary.BigEndian, &dataTaskID)

	return dataTaskID, nil
}

func (rw *ReadWriter) ReadMessageID() (uint64, error) {
	dataByte := make([]byte, rw.MessageIDPlace)
	_, err := rw.Read(dataByte)
	if err != nil {
		return 0, err
	}
	dataBuffer := bytes.NewBuffer(dataByte)
	var dataMesageID uint64
	binary.Read(dataBuffer, binary.BigEndian, &dataMesageID)

	return dataMesageID, nil
}

func (rw *ReadWriter) ReadPackData(length int64) ([]byte, error) {
	dataByte := make([]byte, length)
	_, err := rw.Read(dataByte)
	if err != nil {
		return nil, err
	}
	return dataByte, nil
}

func (rw *ReadWriter) WritePack(taskID uint64, dataByte []byte) (uint64, error) {
	if len(dataByte) == 0 {
		return 0, errors.New("not data")
	}
	rw.MessageID++
	dataLength := uint64(len(dataByte))
	encodeData := dataByte
	if rw.Coding != nil && rw.Coding.Encode != nil {
		encodeData = rw.Coding.Encode(dataByte)
		dataLength = uint64(len(encodeData))
	}
	respPackLen := make([]byte, rw.LenPlace)
	binary.BigEndian.PutUint64(respPackLen, dataLength)

	respPackTaskID := make([]byte, rw.TaskIdPlace)
	binary.BigEndian.PutUint64(respPackTaskID, taskID)

	respMessageID := make([]byte, rw.MessageIDPlace)
	binary.BigEndian.PutUint64(respMessageID, rw.MessageID)

	buffer := &bytes.Buffer{}
	_, err := buffer.Write(respPackLen)
	if err != nil {
		return 0, err
	}

	_, err = buffer.Write(respPackTaskID)
	if err != nil {
		return 0, err
	}

	_, err = buffer.Write(respMessageID)
	if err != nil {
		return 0, err
	}

	_, err = buffer.Write(encodeData)
	if err != nil {
		return 0, err
	}

	_, err = rw.Write(buffer.Bytes())
	if err != nil {
		return 0, err
	}
	return rw.MessageID, rw.Flush()
}

func (rw *ReadWriter) ReadToMessageChan(msgChan chan *Message) (ctx context.Context, cancel context.CancelFunc) {
	ctx, cancel = context.WithCancel(context.Background())
	go func() {
		for {
			msg, err := rw.ReadPack()
			if err != nil {
				cancel()
				return
			}
			msgChan <- msg
		}
	}()
	return ctx, cancel
}
