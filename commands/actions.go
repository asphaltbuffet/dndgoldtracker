package commands

import (
	"log"
	"slices"
	"sort"

	"dndgoldtracker/models"
)

// AddMember creates a new party member in the active member list and gives them last Coin Priority
func AddMember(p *models.Party, name string, xp int, money map[string]int) {
	m := models.Member{Name: name, Level: determineLevel(xp), XP: xp, Coins: money, CoinPriority: len(p.ActiveMembers)}
	p.ActiveMembers = append(p.ActiveMembers, m)
	log.Printf("Welcome to the party %s!\n", m.Name)
}

// ChangeMemberGroup moves a member to a different group e.g. Active to Inactive
func ChangeMemberGroup(srcGroup *[]models.Member, dstGroup *[]models.Member, index int) {
	*dstGroup = append(*dstGroup, (*srcGroup)[index])
	(*dstGroup)[len(*dstGroup)-1].CoinPriority = len(*dstGroup) - 1
	*srcGroup = slices.Delete((*srcGroup), index, index+1)
}

// DistributeCoins distributes coins fairly among party members in a fixed order
// Hands extras out one at a time and rotates coin priority
func DistributeCoins(p *models.Party, money map[string]int) {
	numMembers := len(p.ActiveMembers)
	if numMembers == 0 {
		log.Println("No members to distribute coins to.")
		return
	}

	// Initialize coin maps if not already set
	for i := range p.ActiveMembers {
		if p.ActiveMembers[i].Coins == nil {
			p.ActiveMembers[i].Coins = make(map[string]int)
		}
	}

	// Helper function to distribute a specific coin type
	distributeCoin := func(coinType string, coinAmount int) {
		each := coinAmount / numMembers
		remainder := coinAmount % numMembers

		// Assign evenly to each member
		for i := range p.ActiveMembers {
			log.Printf("Adding %d of %s to %s's wallet\n", each, coinType, p.ActiveMembers[i].Name)
			p.ActiveMembers[i].Coins[coinType] += each
		}

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
	for _, coinType := range models.CoinOrder {
		amount, exists := money[coinType]
		if exists {
			distributeCoin(coinType, amount)
		}
	}
}

// DistributeExperience distributes XP and checks for level-ups
func DistributeExperience(p *models.Party, xp int) {
	for i := range p.ActiveMembers {
		p.ActiveMembers[i].XP += xp / len(p.ActiveMembers)
		checkLevelUp(&p.ActiveMembers[i])
	}

	log.Println("XP added!")
}

// GetFirstCoinPriority returns a member of the party with the lowest priority.
// returns -1 if there are no members
func GetFirstCoinPriority(p *models.Party) int {
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

func checkLevelUp(member *models.Member) {
	for member.Level < len(models.XpThresholds)-1 {
		if member.XP < models.XpThresholds[member.Level] {
			break
		}
		member.Level++
		log.Printf("🎉 %s leveled up to Level %d! 🎉\n", member.Name, member.Level)
	}
}

// Determines the level of a character for a given amount of xp
func determineLevel(xp int) int {
	for i := range models.XpThresholds {
		if xp < models.XpThresholds[i] {
			return i
		}
	}

	// max level
	return 20
}
