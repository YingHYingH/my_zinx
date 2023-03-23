package znet

import (
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
		// 阻塞等待客户端连接，处理客户端连接业务
		for {
			// 如果有客户端连接，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("accept error:", err)
				continue
			}

			// 已经与客户端建立连接，做一些业务操作
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("read error:", err)
						continue
					}
					fmt.Printf("recv client buf %s,cnt %d\n", buf, cnt)

					if _, err := conn.Write(buf[0:cnt]); err != nil {
						fmt.Println("write back buf error:", err)
						continue
					}
				}
			}()
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

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}
