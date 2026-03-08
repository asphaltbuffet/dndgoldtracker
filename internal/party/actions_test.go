package party

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test XP Distribution
func TestDistributeExperience(t *testing.T) {
	party := Party{
		ActiveMembers: []Member{
			{Name: "Keg", Level: 1, XP: 0},
			{Name: "Rowan", Level: 1, XP: 0},
		},
	}

	xpToAdd := 100
	expectedXP := xpToAdd / len(party.ActiveMembers)

	DistributeExperience(&party, xpToAdd)

	// Check if XP was distributed correctly
	for _, member := range party.ActiveMembers {
		if member.XP != expectedXP {
			t.Errorf("Expected %d XP, but got %d for %s", expectedXP, member.XP, member.Name)
		}
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
				{Name: "Keg", CoinPriority: 5, Coins: map[string]int{}},
			},
			want: 0,
		},
		{
			name: "lowest priority is first in slice",
			members: []Member{
				{Name: "Keg", CoinPriority: 0, Coins: map[string]int{}},
				{Name: "Rowan", CoinPriority: 1, Coins: map[string]int{}},
				{Name: "Fred", CoinPriority: 2, Coins: map[string]int{}},
			},
			want: 0,
		},
		{
			name: "lowest priority is last in slice",
			members: []Member{
				{Name: "Keg", CoinPriority: 2, Coins: map[string]int{}},
				{Name: "Rowan", CoinPriority: 1, Coins: map[string]int{}},
				{Name: "Fred", CoinPriority: 0, Coins: map[string]int{}},
			},
			want: 2,
		},
		{
			name: "lowest priority is in the middle",
			members: []Member{
				{Name: "Keg", CoinPriority: 2, Coins: map[string]int{}},
				{Name: "Rowan", CoinPriority: 0, Coins: map[string]int{}},
				{Name: "Fred", CoinPriority: 1, Coins: map[string]int{}},
			},
			want: 1,
		},
		{
			name: "duplicate minimum priorities — returns first occurrence",
			members: []Member{
				{Name: "Keg", CoinPriority: 0, Coins: map[string]int{}},
				{Name: "Rowan", CoinPriority: 0, Coins: map[string]int{}},
				{Name: "Fred", CoinPriority: 1, Coins: map[string]int{}},
			},
			want: 0,
		},
		{
			name: "all members share the same priority",
			members: []Member{
				{Name: "Keg", CoinPriority: 3, Coins: map[string]int{}},
				{Name: "Rowan", CoinPriority: 3, Coins: map[string]int{}},
				{Name: "Fred", CoinPriority: 3, Coins: map[string]int{}},
			},
			want: 0,
		},
		{
			// negative priority values are not prevented by the model
			name: "negative priority values",
			members: []Member{
				{Name: "Keg", CoinPriority: -1, Coins: map[string]int{}},
				{Name: "Rowan", CoinPriority: 0, Coins: map[string]int{}},
				{Name: "Fred", CoinPriority: 1, Coins: map[string]int{}},
			},
			want: 0,
		},
		{
			// duplicate member names — the model has no name uniqueness constraint
			name: "duplicate names — lowest priority is second occurrence",
			members: []Member{
				{Name: "fake", CoinPriority: 1, Coins: map[string]int{}},
				{Name: "fake", CoinPriority: 0, Coins: map[string]int{}},
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
	party := Party{
		ActiveMembers: []Member{
			{Name: "Keg", CoinPriority: 0, Coins: make(map[string]int)},
			{Name: "Rowan", CoinPriority: 1, Coins: make(map[string]int)},
			{Name: "Fred", CoinPriority: 2, Coins: make(map[string]int)},
		},
	}

	// Coins to distribute
	money := map[string]int{
		Platinum: 10,
		Gold:     8,
		Electrum: 7,
		Silver:   5,
		Copper:   3,
	}

	// Call the function to distribute the coins
	DistributeCoins(&party, money)

	// Check the distribution of coins
	tests := []struct {
		memberName       string
		expectedPlatinum int
		expectedGold     int
		expectedElectrum int
		expectedSilver   int
		expectedCopper   int
	}{
		{"Keg" /*PP*/, 4 /*GP*/, 3 /*EP*/, 2 /*SP*/, 2 /*CP*/, 1},
		{"Rowan" /*PP*/, 3 /*GP*/, 2 /*EP*/, 3 /*SP*/, 2 /*CP*/, 1},
		{"Fred" /*PP*/, 3 /*GP*/, 3 /*EP*/, 2 /*SP*/, 1 /*CP*/, 1},
	}

	// Iterate through the test cases and compare expected vs actual
	for _, test := range tests {
		t.Run(test.memberName, func(t *testing.T) {
			member := getMemberByName(party.ActiveMembers, test.memberName)

			if member.Coins[Platinum] != test.expectedPlatinum {
				t.Errorf("%s's platinum: expected %d, got %d", test.memberName, test.expectedPlatinum, member.Coins[Platinum])
			}
			if member.Coins[Gold] != test.expectedGold {
				t.Errorf("%s's gold: expected %d, got %d", test.memberName, test.expectedGold, member.Coins[Gold])
			}
			if member.Coins[Electrum] != test.expectedElectrum {
				t.Errorf("%s's electrum: expected %d, got %d", test.memberName, test.expectedElectrum, member.Coins[Electrum])
			}
			if member.Coins[Silver] != test.expectedSilver {
				t.Errorf("%s's silver: expected %d, got %d", test.memberName, test.expectedSilver, member.Coins[Silver])
			}
			if member.Coins[Copper] != test.expectedCopper {
				t.Errorf("%s's copper: expected %d, got %d", test.memberName, test.expectedCopper, member.Coins[Copper])
			}
		})
	}
}

// Helper function to get PartyMember by name
func getMemberByName(members []Member, name string) *Member {
	for i := range members {
		if members[i].Name == name {
			return &members[i]
		}
	}
	return nil
}
