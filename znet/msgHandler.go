package znet

import (
	"fmt"
	"my_zinx/utils"
	"my_zinx/ziface"
)

type MsgHandler struct {
	// 存放每个MsgID对应的处理方法
	Apis map[uint32]ziface.IRouter
	// 负责Worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	// 业务工作Worker池的worker数量
	WorkerPoolSize uint32
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           map[uint32]ziface.IRouter{},
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
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

// 启动一个Worker工作池（开启工作池的动作只能发生一次，一个zinx框架只能有一个Worker工作池）
func (m *MsgHandler) StartWorkerPool() {
	// 根据workerPoolSize分别开启worker
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		// 当前的worker对应的channel消息队列开辟空间
		m.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go m.StartOneWorker(i, m.TaskQueue[i])
	}
}

// 启动一个Worker工作流程
func (m *MsgHandler) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("workerID = ", workerID, " is started...")

	// 阻塞等待消息队列的消息
	for {
		select {
		case request := <-taskQueue:
			m.DoMsgHandler(request)
		}
	}
}

// 将消息交给TaskQueue，由worker处理
func (m *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	// 将消息平均分配给不同worker
	workID := request.GetConnection().GetConnID() % m.WorkerPoolSize
	// 将消息发送给对应worker的taskQueue
	m.TaskQueue[workID] <- request
}
