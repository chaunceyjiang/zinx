package znet

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
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

//Test Handle
func (b *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	fmt.Println("recv from client : msgId=", request.GetMsgId(), ", data=", string(request.GetData()))
	request.GetConnection().SendMsg(request.GetMsgId()+1, request.GetData())
}



func TestServer_AddRouter(t *testing.T) {
	server := NewServer("test")
	server.AddRouter(&PingRouter{})
	go server.Serve()
	time.Sleep(1 * time.Second)
	addr, _ := net.ResolveTCPAddr(utils.DEFAULT_IP_VERSION, fmt.Sprintf("%s:%d", utils.DEFAULT_IP, utils.DEFAULT_PORT))
	conn, err := net.DialTCP(utils.DEFAULT_IP_VERSION, nil, addr)
	assert.Equal(t, nil, err, "dial tcp")

	dp := NewDataPack()
	msg, _ := dp.Pack(NewMsgPackage(1, []byte("hello world")))

	_, err = conn.Write(msg)
	assert.Equal(t, nil, err, "write tcp")
	head := make([]byte, dp.GetHeadLen())

	io.ReadFull(conn, head)
	m, _ := dp.Unpack(head)
	body := make([]byte, m.GetDataLen())
	io.ReadFull(conn, body)

	assert.Equal(t, true, m.GetMsgId() == 2)
	assert.Equal(t, "hello world", string(body))
}
