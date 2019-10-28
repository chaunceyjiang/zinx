package ziface

// IRequest 包装客户端的请求信息
type IRequest interface {
	// 获取连接信息
	GetConnection() IConnection

	// 获取请求消息的数据
	GetData() []byte
	// 获取请求消息的id
	GetMsgId() uint32
}

