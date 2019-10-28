package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Server struct {
	Name      string     // 服务器名称
	IPVersion string     // ip版本
	IP        string     // 服务器绑定的地址
	Port      int        // 服务器监听的端口
	msgHandle ziface.ImsgHandle //当前Server的消息管理模块，用来绑定MsgId和对应的处理方法
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	//panic("implement me")
	s.msgHandle.AddRouter(msgId,router)
}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server listener at IP: %s, Port %d , is starting...\n", s.IP, s.Port)
	go func() {
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
		var connID uint32 = 0
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("AcceptTcp err", err)
				continue
			}

			dealConn := NewConnection(conn, connID, s.msgHandle)
			go dealConn.Start()
			connID++
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
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.Port,
		msgHandle: NewMsgHandle(),
	}
	return s

}
