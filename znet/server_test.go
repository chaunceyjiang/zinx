package znet

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"zinx/utils"
)

func TestServer(t *testing.T) {
	server := NewServer("test")
	server.Serve()

	addr, _ := net.ResolveTCPAddr(utils.DEFAULT_IP_VERSION, fmt.Sprintf("%s:%d", utils.DEFAULT_IP, utils.DEFAULT_PORT))
	conn, err := net.DialTCP(utils.DEFAULT_IP_VERSION, nil, addr)
	assert.Equal(t, nil, err, "dial tcp")
	_, err = conn.Write([]byte("你好"))
	assert.Equal(t, nil, err, "write tcp")

}
