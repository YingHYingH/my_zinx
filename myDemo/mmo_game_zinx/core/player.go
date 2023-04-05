package core

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"my_zinx/myDemo/mmo_game_zinx/pb"
	"my_zinx/ziface"
	"sync"
)

type Player struct {
	PID  int32
	Conn ziface.IConnection // 当前玩家的连接，用于和客户端连接
	X    float32
	Y    float32 // 高度
	Z    float32
	V    float32 // 旋转的角度 0-360
}

/*
	Player ID生成器
*/
var PIDGen int32 = 1
var PIDLock sync.Mutex

// 创建一个玩家
func NewPlayer(conn ziface.IConnection) *Player {
	// 生成一个玩家ID
	var pID int32
	func() {
		PIDLock.Lock()
		pID = PIDGen
		PIDGen++
		defer PIDLock.Unlock()
	}()
	// 创建一个玩家对象
	p := &Player{
		PID:  pID,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)), // 随机在160坐标点基于X轴偏移
		Y:    0,
		Z:    float32(140 + rand.Intn(20)), // 随机在140坐标点基于Y轴偏移
		V:    0,
	}
	return p
}

/*
	提供一个发送给客户端消息的方法，主要是将pb的protobuf数据序列化之后，再调用zinx的SendMsg方法
*/
func (p *Player) SendMsg(msgID uint32, data proto.Message) {
	// 将message结构体序列化转化成二进制
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err: ", err)
		return
	}
	// 将二进制文件通过zinx的sendMsg将数据发送给客户端
	if p.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}
	if err = p.Conn.SendMsg(msgID, msg); err != nil {
		fmt.Println("Player sendMsg err: ", err)
		return
	}
	return
}

// 告知客户端玩家PID，同步已经生成的玩家ID给客户端
func (p *Player) SyncPID() {
	// 组建MsgID:1的proto数据
	protoMsg := &pb.SyncPID{PID: p.PID}
	p.SendMsg(1, protoMsg)
}

// 广播玩家自己的出生地点
func (p *Player) BroadCastStartPosition() {
	// 组建MsgID:200的proto数据
	protoMsg := &pb.BroadCast{
		PID: p.PID,
		Tp:  2,
		Data: &pb.BroadCast_P{P: &pb.Position{
			X: p.X,
			Y: p.Y,
			Z: p.Z,
			V: p.V,
		}},
	}
	p.SendMsg(200, protoMsg)
}

func (p *Player) Talk(content string) {
	//  组建MsgID200的proto数据
	protoMsg := &pb.BroadCast{
		PID:  p.PID,
		Tp:   1,
		Data: &pb.BroadCast_Content{Content: content},
	}
	// 得到当前所有在线玩家
	players := WorldMgrOjb.GetAllPlayers()
	// 向所有玩家发送消息
	for _, player := range players {
		player.SendMsg(200, protoMsg)
	}
}

// 同步玩家上线的位置消息
func (p *Player) SyncSurrounding() {
	pIDs := WorldMgrOjb.AOIManager.GetPIDsByPos(p.X, p.Z)
	players := make([]*Player, 0, len(pIDs))
	for _, pID := range pIDs {
		players = append(players, WorldMgrOjb.GetPlayerByPID(int32(pID)))
	}
	protoMsg := &pb.BroadCast{
		PID: p.PID,
		Tp:  2,
		Data: &pb.BroadCast_P{P: &pb.Position{
			X: p.X,
			Y: p.Y,
			Z: p.Z,
			V: p.V,
		}},
	}
	for _, player := range players {
		player.SendMsg(200, protoMsg)
	}
	playersProtoMsg := make([]*pb.Player, 0, len(pIDs))
	for _, player := range players {
		playersProtoMsg = append(playersProtoMsg, &pb.Player{
			PID: player.PID,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		})
	}
	syncPlayersProtoMsg := &pb.SyncPlayers{Ps: playersProtoMsg}
	p.SendMsg(202, syncPlayersProtoMsg)
}

func (p *Player) UpdatePos(x, y, z, v float32) {
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v

	protoMsg := &pb.BroadCast{
		PID: p.PID,
		Tp:  4,
		Data: &pb.BroadCast_P{P: &pb.Position{
			X: x,
			Y: y,
			Z: z,
			V: v,
		}},
	}
	players := p.GetSurroundingPlayers()
	for _, player := range players {
		player.SendMsg(200, protoMsg)
	}
}

func (p *Player) GetSurroundingPlayers() []*Player {
	pIDs := WorldMgrOjb.AOIManager.GetPIDsByPos(p.X, p.Z)
	players := make([]*Player, 0, len(pIDs))
	for _, pID := range pIDs {
		players = append(players, WorldMgrOjb.GetPlayerByPID(int32(pID)))
	}
	return players
}
