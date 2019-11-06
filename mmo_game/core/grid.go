package core

import (
	"fmt"
	"sync"
)

type Grid struct {
	GID       int // 格子id
	MinX      int // 格子左边界
	MaxX      int
	MinY      int
	MaxY      int          // 格子下边界
	playerIDs map[int]bool // 当前格子内的晚间或则物体的id
	pIDLock   sync.RWMutex
}

func NewGrid(gid, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:       gid,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[int]bool),
	}
}

// 向当前格子中添加一个玩家
func (g *Grid) Add(playID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()
	g.playerIDs[playID] = true
}

func (g *Grid) Remove(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()
	delete(g.playerIDs, playerID)
}

func (g *Grid) GetPlayerIDs() (playerIds []int) {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()

	for k, _ := range g.playerIDs {
		playerIds = append(playerIds, k)
	}
	return
}

func (g *Grid) String() string {
	return fmt.Sprintf("Grid id: %d,minX: %d, maxX: %d, minY: %d, maxY: %d, playIDs: %v",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}
