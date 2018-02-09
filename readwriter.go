package pnet

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"net"
)

type Coding struct {
	Encode func([]byte) []byte
	Decode func([]byte) ([]byte, error)
}

type ReadWriter struct {
	*bufio.Writer
	*bufio.Reader
	LenPlace    int
	TaskIdPlace int
	Coding      *Coding
}

type Message struct {
	Length  int64
	TaskId  int64
	RawData []byte
	Data    []byte
	err     error
}

func NewReaderWriterFromConn(conn net.Conn, coding *Coding) *ReadWriter {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	rw := &ReadWriter{
		Reader:      reader,
		Writer:      writer,
		LenPlace:    8,
		TaskIdPlace: 8,
		Coding:      coding,
	}
	return rw
}

func (rw *ReadWriter) ReadPack() (*Message, error) {
	dataLength, err := rw.ReadPackLen()
	msg := new(Message)
	if err != nil {
		return nil, err
	}
	dataTaskId, err := rw.ReadPackTaskId()
	if err != nil {
		return nil, err
	}
	data, err := rw.ReadPackData(dataLength)
	if err != nil {
		return nil, err
	}
	msg.Length = dataLength
	msg.TaskId = dataTaskId
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

func (rw *ReadWriter) ReadPackTaskId() (int64, error) {
	dataByte := make([]byte, rw.TaskIdPlace)
	_, err := rw.Read(dataByte)
	if err != nil {
		return 0, err
	}
	dataBuffer := bytes.NewBuffer(dataByte)
	var dataTaskId int64
	binary.Read(dataBuffer, binary.BigEndian, &dataTaskId)

	return dataTaskId, nil
}

func (rw *ReadWriter) ReadPackData(length int64) ([]byte, error) {
	dataByte := make([]byte, length)
	_, err := rw.Read(dataByte)
	if err != nil {
		return nil, err
	}
	return dataByte, nil
}

func (rw *ReadWriter) WritePack(taskId uint64, dataByte []byte) error {
	if len(dataByte) == 0 {
		return errors.New("not data")
	}
	dataLength := uint64(len(dataByte))
	encodeData := dataByte
	if rw.Coding != nil && rw.Coding.Encode != nil {
		encodeData = rw.Coding.Encode(dataByte)
		dataLength = uint64(len(encodeData))
	}
	respPackLen := make([]byte, rw.LenPlace)
	binary.BigEndian.PutUint64(respPackLen, dataLength)

	respPackTaskId := make([]byte, rw.TaskIdPlace)
	binary.BigEndian.PutUint64(respPackTaskId, taskId)

	buffer := &bytes.Buffer{}
	_, err := buffer.Write(respPackLen)
	if err != nil {
		return err
	}

	_, err = buffer.Write(respPackTaskId)
	if err != nil {
		return err
	}

	_, err = buffer.Write(encodeData)
	if err != nil {
		return err
	}

	_, err = rw.Write(buffer.Bytes())
	if err != nil {
		return err
	}
	return rw.Flush()
}
