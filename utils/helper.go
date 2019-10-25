package utils

import (
	"fmt"
	"net"
)

const (
	DEFAULT_IP_VERSION = "tcp4"
	DEFAULT_IP         = "0.0.0.0"
	DEFAULT_PORT       = 7777
)

// DefaultHandFunc 默认的业务处理handler
func DefaultHandFunc(conn *net.TCPConn, data []byte, cnt int) error {
	// 回显
	_, err := conn.Write(data[:cnt])
	if err != nil {
		fmt.Println("HandFunc write err", err)
		return fmt.Errorf("%w HandFunc write err", err)
	}
	return nil
}
