package storage

import (
	"encoding/json"
	"os"

	"dndgoldtracker/internal/party"
)

// SaveParty writes party data to a JSON file
func SaveParty(p *party.Party) error {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("party.json", data, 0o644)
}

// LoadParty loads party data from a JSON file
func LoadParty() (party.Party, error) {
	data, err := os.ReadFile("party.json")
	if err != nil {
		return party.Party{}, err
	}
	var p party.Party
	err = json.Unmarshal(data, &p)
	return p, err
}
