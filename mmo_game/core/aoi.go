package core

import "fmt"

type AOIManager struct {
	MinX  int
	MaxX  int
	CntsX int

	MinY  int
	MaxY  int
	CntsY int

	grids map[int]*Grid
}

func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoi := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		grids: make(map[int]*Grid),
	}
	for y := 0; y < cntsY; y++ {
		//
		for x := 0; x < cntsX; x++ {
			// 计算格子的id
			gid := y*cntsX + x
			// 初始化一个格子
			aoi.grids[gid] = NewGrid(gid,
				aoi.MinX+x*aoi.gridWidth(),
				aoi.MinX+(x+1)*aoi.gridWidth(),
				aoi.MinY+y*aoi.gridLength(),
				aoi.MinY+(y+1)*aoi.gridLength())
		}
	}
	return aoi
}

// 得到每个格子的在x轴的宽度
func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

// 得到每个格子的在x轴的长度
func (m *AOIManager) gridLength() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

// 根据格子的gid得到当前周边的九宫格信息
func (m *AOIManager) GetSurroundGridsByGid(gID int) (grIds []*Grid) {
	if _, ok := m.grids[gID]; !ok {
		return
	}

	// 将当前的gid 天剑到九宫格中
	grIds = append(grIds, m.grids[gID])

	// 格局gid得到当前格子所有的x轴的编号
	idx := gID % m.CntsX

	// 判断当前idx 左边是否还有格子,若idx等于0 表示当前格子靠着左边墙
	if idx > 0 {
		grIds = append(grIds, m.grids[gID-1])
	}
	// 判断当前idx 右边是否还有格子, 若idx 等于 cntsX -1 表示当前格子靠着右边墙
	if idx < m.CntsX-1 {
		grIds = append(grIds, m.grids[gID+1])
	}

	// 将x轴当前的格子取出,计算每个格子上下是否有上下格子,用来统计九宫格

	// 得到当前x轴的格子的Id的集合
	gridsX := make([]int, 0, len(grIds))

	for i, v := range grIds {
		gridsX[i] = v.GID
	}

}

func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManagr:\nminX:%d, maxX:%d, cntsX:%d, minY:%d, maxY:%d, cntsY:%d\n Grids in AOI Manager:\n",
		m.MinX, m.MaxX, m.CntsX, m.MinY, m.MaxY, m.CntsY)
	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}

	return s
}
