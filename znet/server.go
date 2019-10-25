package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Server struct {
	Name      string // 服务器名称
	IPVersion string // ip版本
	IP        string // 服务器绑定的地址
	Port      int
}

func (s *Server) Start() {
	fmt.Printf("[START] Server listener at IP: %s, Port %d , is starting...\n", s.IP, s.Port)
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("resolve tcp addr err ", err)
		return
	}
	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println("listen ", s.IPVersion, addr.String(), "err", err)
	}
	fmt.Println("start Zinx server  ", s.Name, "success, now listening")
	go func() {
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("AcceptTcp err", err)
				continue
			}
			go func() {
				for {
					buf := make([]byte, 512)
					n, err := conn.Read(buf)
					if err != nil {
						fmt.Println("recv buf err", err)
						continue
					}

					if _, err := conn.Write(buf[:n]); err != nil {
						fmt.Println("Write back buf err", err)
						continue
					}
				}
			}()
		}
	}()

}

func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server, name ", s.Name)

}

func (s *Server) Serve() {
	s.Start()

	select {}
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: utils.DEFAULT_IP_VERSION,
		IP:        utils.DEFAULT_IP,
		Port:      utils.DEFAULT_PORT,
	}
	return s

}
