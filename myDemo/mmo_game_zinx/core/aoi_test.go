package core

import (
	"fmt"
	"testing"
)

func TestNewAOIManager(t *testing.T) {
	// 初始化AOIManager
	aoiManager := NewAOIManager(0, 250, 5, 0, 250, 5)
	// 打印AOIManager

	fmt.Println(aoiManager)
}

func TestAOIManagerSurroundGridsByGID(t *testing.T) {
	// 初始化AOIManager
	aoiManager := NewAOIManager(0, 250, 5, 0, 250, 5)

	for gID := range aoiManager.grids {
		// 得到当前gID周边九宫格信息
		grids := aoiManager.GetSurroundGridsByGID(gID)
		fmt.Println("gID: ", gID, " grids len= ", len(grids))
		gIDs := make([]int, 0, len(grids))
		for _, grid := range grids {
			gIDs = append(gIDs, grid.GID)
		}
		fmt.Println("surrounding grid IDs are ", gIDs)
	}
}
