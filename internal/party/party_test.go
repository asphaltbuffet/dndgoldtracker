package party

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test XP Distribution
func TestDistributeExperience(t *testing.T) {
	p := Party{
		ActiveMembers: []Member{
			{Name: "Keg", Level: 1, XP: 0},
			{Name: "Rowan", Level: 1, XP: 0},
		},
	}

	xpToAdd := 100
	expectedXP := xpToAdd / len(p.ActiveMembers)

	p.DistributeExperience(xpToAdd)

	// Check if XP was distributed correctly
	for _, member := range p.ActiveMembers {
		assert.Equal(t, expectedXP, member.XP)
	}
}

func TestGetFirstCoinPriority(t *testing.T) {
	tests := []struct {
		name    string
		members []Member
		want    int
	}{
		{
			name: "single member always returns index 0",
			members: []Member{
				{Name: "Keg", CoinPriority: 5, Coins: map[Coin]int{}},
			},
			want: 0,
		},
		{
			name: "lowest priority is first in slice",
			members: []Member{
				{Name: "Keg", CoinPriority: 0, Coins: map[Coin]int{}},
				{Name: "Rowan", CoinPriority: 1, Coins: map[Coin]int{}},
				{Name: "Fred", CoinPriority: 2, Coins: map[Coin]int{}},
			},
			want: 0,
		},
		{
			name: "lowest priority is last in slice",
			members: []Member{
				{Name: "Keg", CoinPriority: 2, Coins: map[Coin]int{}},
				{Name: "Rowan", CoinPriority: 1, Coins: map[Coin]int{}},
				{Name: "Fred", CoinPriority: 0, Coins: map[Coin]int{}},
			},
			want: 2,
		},
		{
			name: "lowest priority is in the middle",
			members: []Member{
				{Name: "Keg", CoinPriority: 2, Coins: map[Coin]int{}},
				{Name: "Rowan", CoinPriority: 0, Coins: map[Coin]int{}},
				{Name: "Fred", CoinPriority: 1, Coins: map[Coin]int{}},
			},
			want: 1,
		},
		{
			name: "duplicate minimum priorities — returns first occurrence",
			members: []Member{
				{Name: "Keg", CoinPriority: 0, Coins: map[Coin]int{}},
				{Name: "Rowan", CoinPriority: 0, Coins: map[Coin]int{}},
				{Name: "Fred", CoinPriority: 1, Coins: map[Coin]int{}},
			},
			want: 0,
		},
		{
			name: "all members share the same priority",
			members: []Member{
				{Name: "Keg", CoinPriority: 3, Coins: map[Coin]int{}},
				{Name: "Rowan", CoinPriority: 3, Coins: map[Coin]int{}},
				{Name: "Fred", CoinPriority: 3, Coins: map[Coin]int{}},
			},
			want: 0,
		},
		{
			// negative priority values are not prevented by the model
			name: "negative priority values",
			members: []Member{
				{Name: "Keg", CoinPriority: -1, Coins: map[Coin]int{}},
				{Name: "Rowan", CoinPriority: 0, Coins: map[Coin]int{}},
				{Name: "Fred", CoinPriority: 1, Coins: map[Coin]int{}},
			},
			want: 0,
		},
		{
			// duplicate member names — the model has no name uniqueness constraint
			name: "duplicate names — lowest priority is second occurrence",
			members: []Member{
				{Name: "fake", CoinPriority: 1, Coins: map[Coin]int{}},
				{Name: "fake", CoinPriority: 0, Coins: map[Coin]int{}},
			},
			want: 1,
		},
		{
			name:    "empty party",
			members: []Member{},
			want:    -1, // should just use a sentinel value as obviously wrong return
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Party{ActiveMembers: tt.members}

			got := GetFirstCoinPriority(p)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDistributeCoins(t *testing.T) {
	// Create a mock party with 3 members
	mockParty := Party{
		ActiveMembers: []Member{
			{Name: "Keg", CoinPriority: 0, Coins: make(map[Coin]int)},
			{Name: "Rowan", CoinPriority: 1, Coins: make(map[Coin]int)},
			{Name: "Fred", CoinPriority: 2, Coins: make(map[Coin]int)},
		},
	}

	// Coins to distribute
	money := map[Coin]int{
		Platinum: 10,
		Gold:     8,
		Electrum: 7,
		Silver:   5,
		Copper:   3,
	}

	// Call the function to distribute the coins
	DistributeCoins(&mockParty, money)

	// Check coin amounts per member by name.
	// DistributeCoins sorts ActiveMembers in-place, so we look up by name
	// rather than comparing the whole struct (which would also tie us to
	// CoinPriority rotation state and slice ordering).
	wantCoins := map[string]map[Coin]int{
		"Keg":   {Platinum: 4, Gold: 3, Electrum: 2, Silver: 2, Copper: 1},
		"Rowan": {Platinum: 3, Gold: 2, Electrum: 3, Silver: 2, Copper: 1},
		"Fred":  {Platinum: 3, Gold: 3, Electrum: 2, Silver: 1, Copper: 1},
	}

	for _, m := range mockParty.ActiveMembers {
		want, ok := wantCoins[m.Name]
		if !ok {
			t.Errorf("unexpected member %q in result", m.Name)
			continue
		}
		assert.Equal(t, want, m.Coins, "coins for %s", m.Name)
	}
}
