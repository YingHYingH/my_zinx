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
	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	// 创建一个server句柄，使用Zinx的api
	s := znet.NewServer("[zinx v0.8]")

	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})
	// 启动server
	s.Serve()
}

type HelloZinxRouter struct {
	znet.BaseRouter
}

func (p *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	// 先读取客户端的数据，再回写
	err := request.GetConnection().SendMsg(201, []byte("hello zinx..."))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("DoConnectionBegin is Called")
	if err := conn.SendMsg(202, []byte("DoConnection Begin")); err != nil {
		fmt.Println(err)
	}
	// 给当前的连接设置一些属性
	fmt.Println("Set conn Name, Age ...")
	conn.SetProperty("Name", "Dad")
	conn.SetProperty("Age", 18)
}

func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("DoConnectionLost is Called")
	fmt.Println("conn ID = ", conn.GetConnID(), " is Lost")

	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name = ", name)
	}
	if age, err := conn.GetProperty("Age"); err == nil {
		fmt.Println("Age = ", age)
	}
}
