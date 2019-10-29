package ziface

type IConnManager interface {
	// 添加链接
	Add(connection IConnection)
	// 删除链接
	Remove(connection IConnection)

	//获取链接
	Get(connId uint32) (IConnection,error)
	// 删除说有链接
	ClearConn()
	// 当前连接池大小
	Len() int
}