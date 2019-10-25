package ziface

import "net"

type IConnection interface {
	Start() // 启动业务连接 ，让当前连接开始工作

	Stop()  // 停止业务连接 ， 结束当前的连接状态

	GetTCPConnection() *net.TCPConn // 从当前业务连接获取原始的 tcpConn

	GetConnID() uint32 // 获取当前的链接ID

	RemoteAddr() net.Addr // 获取远程客户端地址信息

	Send([]byte) error // 发送数据
}


// HandFunc 定义一个统一处理链接业务的接口
// 这个是所有conn 连接的处理业务的函数接口，
// 第一个参数是原生的socket 第二个参数是客户端数据，第三个参数是数据长度
type HandFunc func(*net.TCPConn, []byte, int) error
