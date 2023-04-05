package apis

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"my_zinx/myDemo/mmo_game_zinx/core"
	"my_zinx/myDemo/mmo_game_zinx/pb"
	"my_zinx/ziface"
	"my_zinx/znet"
)

// 玩家移动
type MoveApi struct {
	znet.BaseRouter
}

func (m *MoveApi) Handle(request ziface.IRequest) {
	// 解析客户端传递的proto协议
	protoMsg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), protoMsg)
	if err != nil {
		fmt.Println("Move: Position Unmarshal err: ", err)
		return
	}
	// 得到当前发送位置是哪个玩家
	pID, err := request.GetConnection().GetProperty("pID")
	if err != nil {
		fmt.Println("GetProperty pID err: ", err)
		return
	}
	// 给其他玩家进行当前玩家的位置信息广播
	player := core.WorldMgrOjb.GetPlayerByPID(pID.(int32))
	player.UpdatePos(protoMsg.X, protoMsg.Y, protoMsg.Z, protoMsg.V)
}
