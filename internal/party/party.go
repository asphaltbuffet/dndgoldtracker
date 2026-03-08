package party

import "fmt"

const (
	// Coin types
	Platinum string = "Platinum"
	Gold     string = "Gold"
	Electrum string = "Electrum"
	Silver   string = "Silver"
	Copper   string = "Copper"
)

var (
	// Define the fixed order of coins
	CoinOrder    = []string{Platinum, Gold, Electrum, Silver, Copper}
	XpThresholds = []int{0, 300, 900, 2700, 6500, 14000, 23000, 34000, 48000, 64000, 85000, 100000, 120000, 140000, 165000, 195000, 225000, 265000, 305000, 355000} // XP values taken for D&D 5e
)

type Member struct {
	Name         string
	Level        int
	XP           int
	Coins        map[string]int
	CoinPriority int
}

type Party struct {
	ActiveMembers   []Member
	InactiveMembers []Member
}

// Display prints the current party state
func (p *Party) Display() {
	fmt.Println("\n=== Party Members ===")
	for _, member := range p.ActiveMembers {
		fmt.Printf("%s (Level %d) - XP: %d, Wallet:%dPP %dGP %dEP %dSP %dCP \n",
			member.Name, member.Level, member.XP,
			member.Coins[Platinum], member.Coins[Gold], member.Coins[Electrum],
			member.Coins[Silver], member.Coins[Copper])
	}
}
