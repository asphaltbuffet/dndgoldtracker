package storage

import (
	"os"
	"path/filepath"
	"testing"

	"dndgoldtracker/internal/party"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// writeTemp writes content to a temp file and returns a cleanup func.
func writeTemp(t *testing.T, content string) string {
	t.Helper()

	f, err := os.CreateTemp(t.TempDir(), "party*.json")
	require.NoError(t, err)

	_, err = f.WriteString(content)
	require.NoError(t, err)
	require.NoError(t, f.Close())

	return f.Name()
}

// loadFrom calls LoadParty against a specific file path by temporarily changing to its directory (LoadParty uses a fixed "party.json" filename).
func loadFrom(t *testing.T, path string) (party.Party, error) {
	t.Helper()

	dir := filepath.Dir(path)
	orig, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(dir))

	t.Cleanup(func() { _ = os.Chdir(orig) })
	require.NoError(t, os.Rename(path, filepath.Join(dir, "party.json")))

	return LoadParty()
}

func TestLoadParty(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		want        party.Party
	}{
		{
			name: "enum value keys",
			fileContent: `{
				"ActiveMembers": [{
				"Name": "Keg", "Level": 1, "XP": 0, "CoinPriority": 0,
				"Coins": {"0": 1, "1": 34, "2": 0, "3": 23, "4": 5}
				}],
				"InactiveMembers": null
			}`,
			want: party.Party{
				ActiveMembers: []party.Member{
					{
						Name: "Keg", Level: 1, XP: 0, CoinPriority: 0,
						Coins: map[party.Coin]int{party.Platinum: 1, party.Gold: 34, party.Electrum: 0, party.Silver: 23, party.Copper: 5},
					},
				},
			},
		},
		{
			name: "string keys",
			fileContent: `{
				"ActiveMembers": [{
				"Name": "Keg", "Level": 1, "XP": 0, "CoinPriority": 0,
				"Coins": {"Platinum": 1, "Gold": 34, "Electrum": 0, "Silver": 23, "Copper": 5}
				}],
				"InactiveMembers": null
			}`,
			want: party.Party{
				ActiveMembers: []party.Member{
					{
						Name: "Keg", Level: 1, XP: 0, CoinPriority: 0,
						Coins: map[party.Coin]int{party.Platinum: 1, party.Gold: 34, party.Electrum: 0, party.Silver: 23, party.Copper: 5},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := writeTemp(t, tt.fileContent)
			p, err := loadFrom(t, path)
			require.NoError(t, err)

			assert.Equal(t, p, tt.want)
		})
	}
}

func TestSaveParty(t *testing.T) {
	dir := t.TempDir()
	orig, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { _ = os.Chdir(orig) })

	p := party.Party{
		ActiveMembers: []party.Member{{
			Name:  "Keg",
			Level: 1,
			Coins: map[party.Coin]int{
				party.Platinum: 1,
				party.Gold:     34,
				party.Electrum: 8,
				party.Silver:   23,
				party.Copper: 69,
			},
		}},
	}

	require.NoError(t, SaveParty(&p))

	data, err := os.ReadFile("party.json")
	require.NoError(t, err)

	got := string(data)

	assert.Contains(t, got, `"Platinum"`)
	assert.NotContains(t, got, `"0"`)
	assert.Contains(t, got, `"Gold"`)
	assert.NotContains(t, got, `"1"`)
	assert.Contains(t, got, `"Electrum"`)
	assert.NotContains(t, got, `"2"`)
	assert.Contains(t, got, `"Silver"`)
	assert.NotContains(t, got, `"3"`)
	assert.Contains(t, got, `"Copper"`)
	assert.NotContains(t, got, `"4"`)
}
