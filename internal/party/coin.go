package party

//go:generate stringer -type=Coin -linecomment
type Coin int

const (
	// Coin types
	Platinum Coin = iota // Platinum
	Gold                 // Gold
	Electrum             // Electrum
	Silver               // Silver
	Copper               // Copper
)

// CoinOrder defines the fixed order of coins
var (
	CoinOrder      = []Coin{Platinum, Gold, Electrum, Silver, Copper}
	CoinOrderNames = []string{Platinum.String(), Gold.String(), Electrum.String(), Silver.String(), Copper.String()}
)

var CoinByName = map[string]Coin{
	Platinum.String(): Platinum,
	Gold.String():     Gold,
	Electrum.String(): Electrum,
	Silver.String():   Silver,
	Copper.String():   Copper,
}
