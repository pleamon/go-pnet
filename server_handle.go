package pnet

type ServerHandler interface {
	// 服务端程序启动成功调用
	Initinize(*Server)
	// 有新客户端连接时调用
	AcceptConnHandle(*Server, *ClientInfo) ([]byte, error)
	// 心跳检测
	HeathHandle(*Server) error
	// 接收到客户端数据时调用
	ReceiveMessageHandle(*Server, *Message) ([]byte, error)
	// 客户端连接断开时调用
	FinishConnHandle(error)
	// 程序错误时调用
	ErrorHandle(int, interface{})
}
