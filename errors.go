package pnet

import "errors"

const (
	// ErrorHandleConn 处理客户端连接错误
	ErrorHandleConn = 1 + iota
)

var (
	// ErrorCheckSumFailed 数据校验失败
	ErrorCheckSumFailed = errors.New("error check sum failed")
)
