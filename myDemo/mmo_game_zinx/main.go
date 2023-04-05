package main

import (
	"fmt"
	"my_zinx/myDemo/mmo_game_zinx/apis"
	"my_zinx/myDemo/mmo_game_zinx/core"
	"my_zinx/myDemo/mmo_game_zinx/pb"
	"my_zinx/ziface"
	"my_zinx/znet"
)

func main() {
	// 创建zinx server句柄
	s := znet.NewServer("MMO Game Zinx")
	// 连接创建和销毁的HOOK函数
	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)
	// 注册路由业务
	s.AddRouter(2, &apis.WorldChatApi{})
	s.AddRouter(3, &apis.MoveApi{})
	// 启动服务
	s.Serve()
}

// 当前客户端建立连接之后的hook函数
func OnConnectionAdd(conn ziface.IConnection) {
	// 创建一个Player对象
	player := core.NewPlayer(conn)
	// 给客户端发送MsgID：1的消息
	player.SyncPID()
	// 给客户端发送MsgID：200的消息
	player.BroadCastStartPosition()
	core.WorldMgrOjb.AddPlayer(player)
	conn.SetProperty("pID", player.PID)

	// 同步周边玩家，告诉他们当前玩家已经上线，广播当前玩家的位置信息
	player.SyncSurrounding()

	fmt.Println("Player PDI = ", player.PID, " online")
}

func OnConnectionLost(conn ziface.IConnection) {
	pID, err := conn.GetProperty("pID")
	if err != nil {
		fmt.Println("GetProperty pID err: ", err)
		return
	}
	// 得到当前玩家周边九宫格内有哪些玩家
	player := core.WorldMgrOjb.GetPlayerByPID(pID.(int32))
	players := player.GetSurroundingPlayers()
	// 给周边玩家广播MsgID201消息
	protoMsg := &pb.SyncPID{
		PID: player.PID,
	}
	for _, p := range players {
		p.SendMsg(201, protoMsg)
	}
	// 将当前玩家从世界管理器删除
	core.WorldMgrOjb.AOIManager.RemoveFromGridByPos(int(player.PID), player.X, player.Z)
	core.WorldMgrOjb.RemovePlayer(player.PID)
}
