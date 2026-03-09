package party

import (
	"fmt"
	"strconv"
)

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

// MarshalText implements encoding.TextMarshaler so map[Coin]int keys serialize
// as coin names (e.g. "Gold") rather than integers in JSON.
func (c Coin) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler. It accepts both the coin
// name ("Gold") and the legacy integer form ("1") so old save files load correctly.
func (c *Coin) UnmarshalText(b []byte) error {
	s := string(b)

	// Try name first (current format)
	if coin, ok := CoinByName[s]; ok {
		*c = coin
		return nil
	}

	// Fall back to integer (legacy format)
	n, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("unknown coin %q", s)
	}
	*c = Coin(n)
	return nil
}
