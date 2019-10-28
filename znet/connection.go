package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	// 当前业务的原生socket
	conn *net.TCPConn
	// 当前连接的ip
	connID uint32

	//当前的连接状态
	isClosed bool

	//该连接的处理方法router
	msgHandle ziface.ImsgHandle

	// 关闭 chan
	done chan struct{}
	//无缓冲管道，用于读、写两个goroutine之间的消息通信
	msgChan chan []byte
}

// StartReader 处理conn读数据的Goroutine
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), "conn reader exit!")
	defer c.Stop()
	dp := NewDataPack()
	for {
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
		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.msgHandle.SendMsgToTaskQueue(req)
		}else {
			go c.msgHandle.DoMsgHandler(req)
		}

		// 调用当前连接的业务(这里执行的是当前conn的绑定的handle方法)
		//if err := c.handlerAPI(c.conn, buf, cnt); err != nil {
		//	fmt.Println("connID ", c.isClosed, "handle is error")
		//	return
		//}
	}

}

func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(),"[ conn writer exit ]")
	for {
		select {
		case msg := <-c.msgChan:
			if _, err := c.conn.Write(msg); err != nil {
				fmt.Println("Write connID id ", c.connID, " error ")
				return
			}
		case <-c.done:
			return
		}
	}
}
// Start 启动业务连接
func (c *Connection) Start() {

	// 开启处理该链接读取到客户端数据之后的请求业务
	go c.StartReader()
	//2 开启用于写回客户端数据流程的Goroutine
	go c.StartWriter()
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
	c.conn.Close()

	c.done <- struct{}{}
	close(c.msgChan)
	close(c.done)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.conn
}

func (c *Connection) GetConnID() uint32 {
	return c.connID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
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

	c.msgChan <- msg
	return nil
}

func NewConnection(conn *net.TCPConn, connIdD uint32, router ziface.ImsgHandle) *Connection {
	c := &Connection{
		conn:      conn,
		connID:    connIdD,
		isClosed:  false,
		done:      make(chan struct{}, 1),
		msgHandle: router,
		msgChan:make(chan []byte),
	}
	return c
}
