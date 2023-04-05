package core

import (
	"fmt"
	"sync"
)

type Grid struct {
	// 格子ID
	GID int
	// 格子左边界
	MinX int
	// 格子右边界
	MaxX int
	// 格子上边界
	MinY int
	// 格子下边界
	MaxY int
	// 当前格子内玩家或者物体成员的ID集合
	playerIDs map[int]bool
	// 保护当前集合的锁
	pIDLock sync.RWMutex
}

func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:       gID,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[int]bool),
	}
}

func (g *Grid) Add(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()
	g.playerIDs[playerID] = true
}

func (g *Grid) Remove(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()
	delete(g.playerIDs, playerID)
}

func (g *Grid) GetPlayerIDs() (playerIDs []int) {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()
	for k := range g.playerIDs {
		playerIDs = append(playerIDs, k)
	}
	return
}

func (g *Grid) String() string {
	return fmt.Sprintf("Grid id: %d, minX: %d, maxX: %d, minY: %d, maxY: %d, playerIDs: %v", g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}
