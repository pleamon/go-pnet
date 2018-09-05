package pnet

import "errors"

const (
	ErrorHandleConn = 1 + iota
)

var (
	ErrorCheckSumFailed = errors.New("error check sum failed")
)
