package party

import (
	"log/slog"
	"slices"
	"sort"
)

type Member struct {
	Name         string
	Level        int
	XP           int
	Coins        map[Coin]int
	CoinPriority int
}

// ChangeMemberGroup moves a member to a different group e.g. Active to Inactive
func ChangeMemberGroup(srcGroup *[]Member, dstGroup *[]Member, index int) {
	*dstGroup = append(*dstGroup, (*srcGroup)[index])
	(*dstGroup)[len(*dstGroup)-1].CoinPriority = len(*dstGroup) - 1
	*srcGroup = slices.Delete((*srcGroup), index, index+1)
}

func (m *Member) UpdateLevel() {
	prev := m.Level
	newLevel := sort.SearchInts(XpThresholds, m.XP+1)

	if newLevel > m.Level {
		m.Level = newLevel
		slog.Info("level up", "member", m.Name, "level", m.Level, "from", prev)
	}
}
