package main

import (
	"fmt"
	"my_zinx/ziface"
	"my_zinx/znet"
)

/*
	基于Zinx框架开发的服务器端应用程序
*/

// ping test自定义路由
type PingRouter struct {
	znet.BaseRouter
}

func (p *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("callback before ping error: ", err)
		return
	}
}

func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping...\n"))
	if err != nil {
		fmt.Println("callback ping error: ", err)
		return
	}
}

func (p *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
	if err != nil {
		fmt.Println("callback after ping error: ", err)
		return
	}
}

func main() {
	// 创建一个server句柄，使用Zinx的api
	s := znet.NewServer("[zinx v0.2]")
	s.AddRouter(&PingRouter{})
	// 启动server
	s.Serve()
}
