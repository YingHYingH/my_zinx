package znet

import (
	"errors"
	"fmt"
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
	// 当前的server添加一个router
	Router ziface.IRouter
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server Listener at IP:%s, Port:%d, is starting\n", s.IP, s.Port)

	go func() {
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

			delConn := NewConnection(conn, cid, s.Router)
			cid++
			go delConn.Start()
		}
	}()
}

func (s *Server) Stop() {

}

func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()

	// TODO做一些启动服务器之后的额外业务

	// 阻塞
	select {}
}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add Router success")
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
		Router:    nil,
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
