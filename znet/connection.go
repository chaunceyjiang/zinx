package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/ziface"
)

type Connection struct {
	// 当前业务的原生socket
	Conn *net.TCPConn
	// 当前连接的ip
	ConnID uint32

	//当前的连接状态
	isClosed bool

	//该连接的处理方法router
	Router ziface.IRouter

	// 关闭 chan
	done chan struct{}
}

// StartReader 处理conn读数据的Goroutine
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), "conn reader exit!")
	defer c.Stop()

	for {
		dp := NewDataPack()
		// 读取客户端发送的数据
		head := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), head); err != nil {
			fmt.Println("recv msg head error", err)
			return
		}
		msg, err := dp.Unpack(head)
		if err != nil {
			fmt.Println("unpack error ", err)
			return
		}
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				return
			}
		}

		msg.SetData(data)
		req := &Request{
			conn: c,
			msg:  msg,
		}

		go func(request ziface.IRequest) {
			//从路由Routers 中找到注册绑定Conn的对应Handle
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(req)
		// 调用当前连接的业务(这里执行的是当前conn的绑定的handle方法)
		//if err := c.handlerAPI(c.Conn, buf, cnt); err != nil {
		//	fmt.Println("connID ", c.isClosed, "handle is error")
		//	return
		//}
	}

}

// Start 启动业务连接
func (c *Connection) Start() {

	// 开启处理该链接读取到客户端数据之后的请求业务
	go c.StartReader()

	for {
		select {
		case <-c.done:
			return
		}
	}
}

// Stop 关闭当前业务连接
func (c *Connection) Stop() {
	if c.isClosed {
		// 已经关闭
		return
	}
	// 关闭当前业务连接
	c.isClosed = true

	// 关闭原生连接
	c.Conn.Close()

	c.done <- struct{}{}

	close(c.done)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when send msg\n")
	}
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg ")
	}
	if _, err := c.Conn.Write(msg); err != nil {
		fmt.Println("Write msg id ", msgId, " error ")
		return errors.New("conn Write error")
	}

	return nil
}

func NewConnection(conn *net.TCPConn, connIdD uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connIdD,
		isClosed: false,
		done:     make(chan struct{}, 1),
		Router:   router,
	}
	return c
}
