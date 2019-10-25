package znet

import "zinx/ziface"

type Request struct {
	// 已经和客户端建立好的业务链接
	conn ziface.IConnection
	// 客户端的请求数据
	data []byte
}
// GetConnection 获取链接信息
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}
// GetData 获取客户端的请求数据
func (r *Request) GetData() []byte {
	return r.data
}



