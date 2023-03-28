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

func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	// 先读取客户端的数据，再回写
	err := request.GetConnection().SendMsg(1, request.GetData())
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	// 创建一个server句柄，使用Zinx的api
	s := znet.NewServer("[zinx v0.5]")
	s.AddRouter(&PingRouter{})
	// 启动server
	s.Serve()
}
