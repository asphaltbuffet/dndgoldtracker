package ui

import (
	"testing"

	"dndgoldtracker/internal/party"

	"charm.land/bubbles/v2/textinput"
	"github.com/stretchr/testify/assert"
)

func TestCheckbox(t *testing.T) {
	unchecked := checkbox("Distribute Money", false)
	assert.Equal(t, "[ ] Distribute Money", unchecked)

	checked := checkbox("Distribute Money", true)
	assert.Contains(t, checked, "[x] Distribute Money")
}

func TestMembersToRows(t *testing.T) {
	t.Run("empty slice returns empty rows", func(t *testing.T) {
		rows := membersToRows([]party.Member{})
		assert.Empty(t, rows)
	})

	t.Run("one member per row", func(t *testing.T) {
		members := []party.Member{
			{Name: "Keg", Level: 3, XP: 900, Coins: map[string]int{
				party.Platinum: 1,
				party.Gold:     2,
				party.Electrum: 3,
				party.Silver:   4,
				party.Copper:   5,
			}},
			{Name: "Rowan", Level: 1, XP: 0, Coins: map[string]int{}},
		}
		rows := membersToRows(members)
		assert.Len(t, rows, 2)
	})

	t.Run("row columns match member data", func(t *testing.T) {
		member := party.Member{
			Name:  "Keg",
			Level: 3,
			XP:    900,
			Coins: map[string]int{
				party.Platinum: 1,
				party.Gold:     2,
				party.Electrum: 3,
				party.Silver:   4,
				party.Copper:   5,
			},
		}
		rows := membersToRows([]party.Member{member})
		assert.Len(t, rows, 1)
		row := rows[0]
		assert.Equal(t, "Keg", row[0], "Name field")
		assert.Equal(t, "900", row[1], "XP field")
		assert.Equal(t, "3", row[2], "Level field")
		assert.Equal(t, "1", row[3], "Platinum")
		assert.Equal(t, "2", row[4], "Gold")
		assert.Equal(t, "3", row[5], "Electrum")
		assert.Equal(t, "4", row[6], "Silver")
		assert.Equal(t, "5", row[7], "Copper")
	})

	t.Run("zero-value coins render as 0", func(t *testing.T) {
		member := party.Member{Name: "Fred", Coins: map[string]int{}}
		rows := membersToRows([]party.Member{member})
		assert.Equal(t, "0", rows[0][3], "Platinum default of zero")
	})
}

func makeInputs(t *testing.T, n int) []textinput.Model {
	t.Helper()

	inputs := make([]textinput.Model, n)
	for i := range inputs {
		inputs[i] = textinput.New()
	}

	return inputs
}

func TestUpdateFocusIndex(t *testing.T) {
	t.Run("wraps forward past submit button back to 0", func(t *testing.T) {
		inputs := makeInputs(t, 3)
		idx := len(inputs) + 1 // one past the submit button
		updateFocusIndex(&idx, inputs)
		assert.Equal(t, 0, idx)
	})

	t.Run("wraps backward below 0 to len(inputs)", func(t *testing.T) {
		inputs := makeInputs(t, 3)
		idx := -1
		updateFocusIndex(&idx, inputs)
		assert.Equal(t, len(inputs), idx)
	})

	t.Run("valid mid-range index stays unchanged", func(t *testing.T) {
		inputs := makeInputs(t, 5)
		idx := 2
		updateFocusIndex(&idx, inputs)
		assert.Equal(t, 2, idx)
	})

	t.Run("index at submit button (len) stays unchanged", func(t *testing.T) {
		inputs := makeInputs(t, 3)
		idx := len(inputs) // valid: points at submit button
		updateFocusIndex(&idx, inputs)
		assert.Equal(t, len(inputs), idx)
	})
}

func TestBuildInputList(t *testing.T) {
	t.Run("contains Submit button text", func(t *testing.T) {
		inputs := makeInputs(t, 2)
		out := buildInputList(inputs, 0, false)
		assert.Contains(t, out, "Submit")
	})

	t.Run("shows 'real' cursor mode when virtualCursor is false", func(t *testing.T) {
		inputs := makeInputs(t, 1)
		out := buildInputList(inputs, 0, false)
		assert.Contains(t, out, "real")
	})

	t.Run("shows 'virtual' cursor mode when virtualCursor is true", func(t *testing.T) {
		inputs := makeInputs(t, 1)
		out := buildInputList(inputs, 0, true)
		assert.Contains(t, out, "virtual")
	})

	t.Run("contains cursor mode help text", func(t *testing.T) {
		inputs := makeInputs(t, 1)
		out := buildInputList(inputs, 0, false)
		assert.Contains(t, out, "ctrl+r to change style")
	})
}
