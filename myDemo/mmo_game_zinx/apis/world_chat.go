package apis

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"my_zinx/myDemo/mmo_game_zinx/core"
	"my_zinx/myDemo/mmo_game_zinx/pb"
	"my_zinx/ziface"
	"my_zinx/znet"
)

// 世界聊天 路由业务
type WorldChatApi struct {
	znet.BaseRouter
}

func (wc *WorldChatApi) Handle(request ziface.IRequest) {
	// 解析客户端传递的proto协议
	protoMsg := &pb.Talk{}
	err := proto.Unmarshal(request.GetData(), protoMsg)
	if err != nil {
		fmt.Println("Talk Unmarshal err: ", err)
		return
	}
	// 当前聊天数据是哪个pID发送的
	pID, err := request.GetConnection().GetProperty("pID")
	if err != nil {
		fmt.Println("request GetProperty err: ", err)
		return
	}
	// 根据pID得到当前玩家player对象
	player := core.WorldMgrOjb.GetPlayerByPID(pID.(int32))
	player.Talk(protoMsg.Content)
}
