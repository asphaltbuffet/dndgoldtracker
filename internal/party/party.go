package party

import (
	"fmt"
	"log/slog"
	"sort"
)

type Party struct {
	ActiveMembers   []Member
	InactiveMembers []Member
}

// AddMember creates a new party member in the active member list and gives them last Coin Priority
func (p *Party) AddMember(name string, xp int, money map[Coin]int) {
	m := Member{
		Name:         name,
		Level:        0,
		XP:           xp,
		Coins:        money,
		CoinPriority: len(p.ActiveMembers),
	}

	slog.Info("new party member", "name", m.Name)

	p.ActiveMembers = append(p.ActiveMembers, m)
}

// DistributeCoins distributes coins fairly among party members in a fixed order
// Hands extras out one at a time and rotates coin priority
func DistributeCoins(p *Party, money map[Coin]int) {
	numMembers := len(p.ActiveMembers)
	if numMembers == 0 {
		slog.Warn("did not distribute money", "party", p, "money", money)
		return
	}

	// Initialize coin maps if not already set
	for i := range p.ActiveMembers {
		if p.ActiveMembers[i].Coins == nil {
			p.ActiveMembers[i].Coins = make(map[Coin]int)
		}
	}

	// Helper function to distribute a specific coin type
	distributeCoin := func(coinType Coin, coinAmount int) {
		each := coinAmount / numMembers
		remainder := coinAmount % numMembers

		// Assign evenly to each member
		for i := range p.ActiveMembers {
			slog.Debug("adding coins", "amount", each, "type", coinType, "member", p.ActiveMembers[i].Name)
			p.ActiveMembers[i].Coins[coinType] += each
		}

		// TODO: This shouldn't sort member order in place. Make a copy first.
		// Sort members by priority for distributing the remainder
		sort.Slice(p.ActiveMembers, func(i, j int) bool {
			return p.ActiveMembers[i].CoinPriority < p.ActiveMembers[j].CoinPriority
		})

		// Distribute excess coins based on priority
		for i := range remainder {
			p.ActiveMembers[i].Coins[coinType]++
		}

		// Rotate priority to balance future distributions
		for i := range p.ActiveMembers {
			p.ActiveMembers[i].CoinPriority = (p.ActiveMembers[i].CoinPriority + 1) % numMembers
		}
	}

	// Distribute coins in the predefined order
	for _, coinType := range CoinOrder {
		amount, exists := money[coinType]
		slog.Info("distributing coin", "type", coinType, "amount", amount)

		if exists {
			distributeCoin(coinType, amount)
		}
	}
}

// GetFirstCoinPriority returns a member of the party with the lowest priority.
// returns -1 if there are no members
func GetFirstCoinPriority(p *Party) int {
	if len(p.ActiveMembers) == 0 {
		return -1
	}
	minIdx := 0

	for i := range p.ActiveMembers {
		if p.ActiveMembers[i].CoinPriority < p.ActiveMembers[minIdx].CoinPriority {
			minIdx = i
		}
	}

	return minIdx
}
