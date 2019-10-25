package znet

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
	"zinx/utils"
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
