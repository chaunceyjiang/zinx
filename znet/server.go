package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Server struct {
	Name        string            // 服务器名称
	IPVersion   string            // ip版本
	IP          string            // 服务器绑定的地址
	Port        int               // 服务器监听的端口
	msgHandle   ziface.ImsgHandle //当前Server的消息管理模块，用来绑定MsgId和对应的处理方法
	connManager ziface.IConnManager

	//该Server的连接创建时Hook函数
	onConnStart    func(conn ziface.IConnection)
	//该Server的连接断开时的Hook函数
	onConnStop func(conn ziface.IConnection)
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	//panic("implement me")
	s.msgHandle.AddRouter(msgId, router)
}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server listener at IP: %s, Port %d , is starting...\n", s.IP, s.Port)
	s.msgHandle.StartWorkerPool() // 启动工作池
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
			// 当前的连接数超过最大连接数,这关闭这个新连接
			if s.connManager.Len() > utils.GlobalObject.MaxConn {
				conn.Close()
				continue
			}
			dealConn := NewConnection(s, conn, connID, s.msgHandle)
			go dealConn.Start()
			connID++
		}
	}()

}

func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server, name ", s.Name)
	// 服务关闭,则清除全部连接

	s.GetConnManager().ClearConn()
}

func (s *Server) Serve() {
	s.Start()

	select {}
}

func (s *Server) GetConnManager() ziface.IConnManager {
	return s.connManager
}
//设置该Server的连接时的Hook函数
func (s *Server) SetOnConnStart(hook func(connection ziface.IConnection)) {
	s.onConnStart = hook
}

//设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hook func(connection ziface.IConnection)) {
	s.onConnStop = hook
}

//调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(connection ziface.IConnection){
	if s.onConnStart != nil {
		fmt.Println("---> CallOnConnStart....")
		s.onConnStart(connection)
	}
}

//调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(connection ziface.IConnection) {
	if s.onConnStop != nil {
		fmt.Println("---> CallOnConnStart....")
		s.onConnStop(connection)
	}
}

func NewServer(name string) ziface.IServer {
	utils.GlobalObject.Reload()
	s := &Server{
		Name:        name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.Port,
		msgHandle:   NewMsgHandle(),
		connManager: NewConnManager(),
	}
	return s

}
