package core

import "fmt"

const (
	AOI_MIN_X  int = 85
	AOI_MAX_X  int = 410
	AOI_CNTS_X int = 10
	AOI_MIN_Y  int = 75
	AOI_MAX_Y  int = 400
	AOI_CNTS_Y int = 20
)

/*
	AOI区域管理模块
*/

type AOIManager struct {
	// 区域左边界
	MinX int
	// 区域右边界
	MaxX int
	// X方向格子的数量
	CntX int
	// 区域上边界
	MinY int
	// 区域下边界
	MaxY int
	// Y方向格子的数量
	CntY int
	// 当前区域中有哪些格子 key 格子ID value 格子对象
	grids map[int]*Grid
}

func NewAOIManager(minX, maxX, cntX, minY, maxY, cntY int) *AOIManager {
	aoiManager := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntX:  cntX,
		MinY:  minY,
		MaxY:  maxY,
		CntY:  cntY,
		grids: make(map[int]*Grid),
	}

	// 给AOI初始化区域的格子进行编号和初始化
	for y := 0; y < cntY; y++ {
		for x := 0; x < cntX; x++ {
			// 计算格子ID
			gID := y*cntX + x
			// 初始化gID格子
			aoiManager.grids[gID] = NewGrid(gID,
				aoiManager.MinX+x*aoiManager.gridWidth(),
				aoiManager.MinX+(x+1)*aoiManager.gridWidth(),
				aoiManager.MinY+y*aoiManager.gridLength(),
				aoiManager.MinY+(y+1)*aoiManager.gridLength())
		}
	}
	return aoiManager
}

func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntX
}

func (m *AOIManager) gridLength() int {
	return (m.MaxY - m.MinY) / m.CntY
}

func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager:\n MinX:%d, MaxX:%d, cntX:%d, minY:%d, maxY:%d, cntY:%d\n Grids in AOIManager:\n", m.MinX, m.MaxX, m.CntX, m.MinY, m.MaxY, m.CntY)
	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}
	return s
}

// 根据格子GID得到周边九宫格格子的ID集合
func (m *AOIManager) GetSurroundGridsByGID(gID int) (grids []*Grid) {
	// 判断gID是否在AOIManager中
	if _, ok := m.grids[gID]; !ok {
		return
	}
	grids = append(grids, m.grids[gID])
	idx := gID % m.CntX
	if idx > 0 {
		grids = append(grids, m.grids[gID-1])
	}
	if idx < m.CntX-1 {
		grids = append(grids, m.grids[gID+1])
	}
	gIDsX := make([]int, 0, len(grids))
	for _, grid := range grids {
		gIDsX = append(gIDsX, grid.GID)
	}
	for _, gID = range gIDsX {
		idy := gID / m.CntY
		if idy > 0 {
			grids = append(grids, m.grids[gID-m.CntX])
		}
		if idy < m.CntY-1 {
			grids = append(grids, m.grids[gID+m.CntX])
		}
	}
	return grids
}

func (m *AOIManager) GetGIDByPos(x, y float32) int {
	idx := (int(x) - m.MinX) / m.gridWidth()
	idy := (int(y) - m.MinY) / m.gridLength()
	return idy*m.CntX + idx
}

// 通过横纵坐标得到周边九宫格内全部的playerIDs
func (m *AOIManager) GetPIDsByPos(x, y float32) (playerIDs []int) {
	// 得到当前玩家的gID
	gID := m.GetGIDByPos(x, y)
	// 通过gID得到周边九宫格信息
	grids := m.GetSurroundGridsByGID(gID)

	for _, grid := range grids {
		playerIDs = append(playerIDs, grid.GetPlayerIDs()...)
	}
	return playerIDs
}

// 添加一个Player到一个格子中
func (m *AOIManager) AddPIDToGrid(pID, gID int) {
	m.grids[gID].Add(pID)
}

// 移除一个格子中的player
func (m *AOIManager) RemovePIDFromGrid(pID, gID int) {
	m.grids[gID].Remove(pID)
}

// 通过gID获取全部的player
func (m *AOIManager) GetPIDsByGID(gID int) (playerIDs []int) {
	return m.grids[gID].GetPlayerIDs()
}

// 通过坐标将player添加到一个格子中
func (m *AOIManager) AddToGridByPos(pID int, x, y float32) {
	gID := m.GetGIDByPos(x, y)
	m.AddPIDToGrid(pID, gID)
}

// 通过坐标把一个player从一个格子中删除
func (m *AOIManager) RemoveFromGridByPos(pID int, x, y float32) {
	gID := m.GetGIDByPos(x, y)
	m.RemovePIDFromGrid(pID, gID)
}
