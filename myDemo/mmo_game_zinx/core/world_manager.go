package core

import "sync"

/*
	当前游戏的世界总管理模块
*/

type WorldManager struct {
	// 当前世界地图AOI的管理模块
	AOIManager *AOIManager
	// 当前全部在线的players集合
	Players map[int32]*Player
	// 保护Players集合的锁
	pLock sync.RWMutex
}

var WorldMgrOjb *WorldManager

// 初始化
func init() {
	WorldMgrOjb = &WorldManager{
		AOIManager: NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNTS_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_Y),
		Players:    make(map[int32]*Player),
	}
}

func (wm *WorldManager) AddPlayer(player *Player) {
	func() {
		wm.pLock.Lock()
		defer wm.pLock.Unlock()
		wm.Players[player.PID] = player
	}()
	// 将player添加到AOIManager
	wm.AOIManager.AddToGridByPos(int(player.PID), player.X, player.Z)
}

func (wm *WorldManager) RemovePlayer(pID int32) {
	player := wm.Players[pID]
	wm.AOIManager.RemoveFromGridByPos(int(pID), player.Y, player.Z)
	func() {
		wm.pLock.Lock()
		defer wm.pLock.Unlock()
		delete(wm.Players, pID)
	}()
}

func (wm *WorldManager) GetPlayerByPID(pID int32) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()
	return wm.Players[pID]
}

func (wm *WorldManager) GetAllPlayers() []*Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()
	players := make([]*Player, 0)
	for _, player := range wm.Players {
		players = append(players, player)
	}
	return players
}
