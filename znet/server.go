package znet

import (
	"errors"
	"fmt"
	"my_zinx/utils"
	"my_zinx/ziface"
	"net"
)

type Server struct {
	// 服务器名称
	Name string
	// 服务器IP版本
	IPVersion string
	// 服务器监听的IP
	IP string
	// 服务器监听的端口
	Port int
	// 当前的server消息管理模块，绑定msgID和处理业务API关系
	MsgHandler ziface.IMsgHandler
	// 该server的连接管理器
	ConnManager ziface.IConnManager
	// 创建连接之后的Hook方法
	OnConnStart func(conn ziface.IConnection)
	// 销毁连接之前的Hook方法
	OnConnStop func(conn ziface.IConnection)
}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name: %s, listenner at IP: %s, Port: %d is starting\n", utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)

	go func() {
		// 开启消息队列和worker工作池
		s.MsgHandler.StartWorkerPool()
		// 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
			return
		}
		// 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen tcp error:", err)
			return
		}
		fmt.Println("start Zinx server success, ", s.Name, ", Listening..")
		var cid uint32
		// 阻塞等待客户端连接，处理客户端连接业务
		for {
			// 如果有客户端连接，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("accept error:", err)
				continue
			}
			// 设置最大连接个数的判断，如果超过最大连接，则关闭新连接
			if s.ConnManager.Len() >= utils.GlobalObject.MaxConn {
				fmt.Println("too many conn")
				conn.Close()
				continue
			}
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] ZInx server name ", s.Name)
	s.ConnManager.ClearConn()
}

func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()

	// TODO做一些启动服务器之后的额外业务

	// 阻塞
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router success")
}

func (s *Server) GetConnManager() ziface.IConnManager {
	return s.ConnManager
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:        utils.GlobalObject.Name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.TcpPort,
		MsgHandler:  NewMsgHandler(),
		ConnManager: NewConnManager(),
	}
	return s
}

// 定义当前客户端连接callback api，以后应该由应用自定义实现
func CallbackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Println("[Conn Handle] CallbackToClient...")
	_, err := conn.Write(data[:cnt])
	if err != nil {
		fmt.Println("write back buf error: ", err)
		return errors.New("CallbackToClient error")
	}
	return nil
}

func (s *Server) SetOnConnStart(f func(conn ziface.IConnection)) {
	s.OnConnStart = f
}

func (s *Server) SetOnConnStop(f func(conn ziface.IConnection)) {
	s.OnConnStop = f
}

func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("call OnConnStart")
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("call OnConnStop")
		s.OnConnStop(conn)
	}
}
