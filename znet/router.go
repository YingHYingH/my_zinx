package znet

import (
	"my_zinx/ziface"
)

// 实现Router时，先嵌入这个BaseRouter基类，然后根据需要对这个基类的方法重写
type BaseRouter struct{}

func (b *BaseRouter) PreHandle(request ziface.IRequest) {
}

func (b *BaseRouter) Handle(request ziface.IRequest) {
}

func (b *BaseRouter) PostHandle(request ziface.IRequest) {
}
