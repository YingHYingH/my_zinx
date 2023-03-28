package znet

import (
	"fmt"
	"my_zinx/ziface"
)

type MsgHandler struct {
	// 存放每个MsgID对应的处理方法
	Apis map[uint32]ziface.IRouter
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{Apis: map[uint32]ziface.IRouter{}}
}

func (m *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	handler, ok := m.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), " is not found! need register")
		return
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (m *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	if _, ok := m.Apis[msgID]; ok {
		// ID已经存在
		return
	}
	m.Apis[msgID] = router
}
