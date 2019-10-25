package znet

import (
	"fmt"
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

	// 该连接的处理方法func
	handlerAPI ziface.HandFunc

	// 关闭 chan
	done     chan struct{}
	readDone chan struct{}
}

// StartReader 处理conn读数据的Goroutine
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), "conn reader exit!")
	defer c.Stop()

	for {
		select {
		case <-c.readDone:
			return
		default:
			// 读取客户端发送的数据
			buf := make([]byte, 512)
			cnt, err := c.Conn.Read(buf)
			if err != nil {
				fmt.Println("recv buf err", err)
				return
			}
			// 调用当前连接的业务(这里执行的是当前conn的绑定的handle方法)
			if err := c.handlerAPI(c.Conn, buf, cnt); err != nil {
				fmt.Println("connID ", c.isClosed, "handle is error")
				return
			}
		}
	}

}

// Start 启动业务连接
func (c *Connection) Start() {

	// 开启处理该链接读取到客户端数据之后的请求业务
	go c.StartReader()

	for {
		select {
		case <-c.done:
			c.readDone <- struct{}{}
			close(c.readDone)
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

func (c *Connection) Send([]byte) error {
	panic("implement me")
}

func NewConnection(conn *net.TCPConn, connIdD uint32, callbackFunc ziface.HandFunc) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connIdD,
		isClosed:   false,
		handlerAPI: callbackFunc,
		done:       make(chan struct{}, 1),
		readDone:   make(chan struct{}, 1),
	}
	return c
}
