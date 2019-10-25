package znet

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
	"zinx/utils"
	"zinx/ziface"
)

func TestServer(t *testing.T) {
	server := NewServer("test")
	go server.Serve()
	time.Sleep(1 * time.Second)
	for i := 0; i < 10; i++ {
		addr, _ := net.ResolveTCPAddr(utils.DEFAULT_IP_VERSION, fmt.Sprintf("%s:%d", utils.DEFAULT_IP, utils.DEFAULT_PORT))
		conn, err := net.DialTCP(utils.DEFAULT_IP_VERSION, nil, addr)
		assert.Equal(t, nil, err, "dial tcp")
		_, err = conn.Write([]byte("hello world"))
		assert.Equal(t, nil, err, "write tcp")
		buf := make([]byte, 512)
		n, err := conn.Read(buf)
		assert.Equal(t, "hello world", string(buf[:n]), string(buf[:n]))
	}

}

type PingRouter struct {
	BaseRouter
}

func (b *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call PingRouter PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping ....\n"))
	if err != nil {
		fmt.Println("call back ping error")
	}
}

//Test Handle
func (b *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

//Test PostHandle
func (b *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call PingRouter PostHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("After ping .....\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

func TestServer_AddRouter(t *testing.T) {
	server := NewServer("test")
	server.AddRouter(&PingRouter{})
	go server.Serve()
	time.Sleep(1 * time.Second)
	addr, _ := net.ResolveTCPAddr(utils.DEFAULT_IP_VERSION, fmt.Sprintf("%s:%d", utils.DEFAULT_IP, utils.DEFAULT_PORT))
	conn, err := net.DialTCP(utils.DEFAULT_IP_VERSION, nil, addr)
	assert.Equal(t, nil, err, "dial tcp")
	_, err = conn.Write([]byte("hello world"))
	assert.Equal(t, nil, err, "write tcp")

	// TODO 这里的测试代码

	//data, _ := ioutil.ReadAll(conn)
	//fmt.Println("read all", string(data))
	//
	////assert.Equal(t, "After ping .....\n", string(data), string(data))
	//time.Sleep(2 * time.Second)
}
