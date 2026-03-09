package ui

import (
	"regexp"
	"testing"

	"dndgoldtracker/internal/party"

	"github.com/stretchr/testify/assert"
)

var ansiEscape = regexp.MustCompile(`\x1B(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])`)

func stripANSI(s string) string {
	return ansiEscape.ReplaceAllString(s, "")
}

// testModel builds a minimal model with one active member, suitable for views
// that call commands.GetFirstCoinPriority.
func testModel(t *testing.T) model {
	t.Helper()

	member := party.Member{
		Name:         "Keg",
		Level:        1,
		XP:           0,
		CoinPriority: 0,
		Coins:        map[party.Coin]int{},
	}
	p := party.Party{
		ActiveMembers: []party.Member{member},
	}

	coinInputs := configureInputs(party.CoinOrderNames)
	xpInputs := configureInputs([]string{xp})
	memberFields := []string{name, xp}
	memberFields = append(memberFields, party.CoinOrderNames...)
	memberInputs := configureInputs(memberFields)

	return model{
		party:               p,
		activeMemberTable:   configureTable(p.ActiveMembers),
		inactiveMemberTable: configureTable(p.InactiveMembers),
		coinInputs:          coinInputs,
		xpInputs:            xpInputs,
		memberInputs:        memberInputs,
		virtualCursor:       true,
	}
}

func TestChoicesView(t *testing.T) {
	m := testModel(t)
	out := stripANSI(choicesView(m))

	t.Run("contains all menu options", func(t *testing.T) {
		assert.Contains(t, out, "Distribute Money")
		assert.Contains(t, out, "Distribute Experience")
		assert.Contains(t, out, "Add Member")
		assert.Contains(t, out, "Activate")
	})

	t.Run("shows coin priority member name", func(t *testing.T) {
		assert.Contains(t, out, "Current Coin Priority is to Keg")
	})

	t.Run("first choice selected shows checked, others unchecked", func(t *testing.T) {
		m.choice = 0
		newOut := choicesView(m)
		// choice 0 is checked, choice 1 is not
		assert.Contains(t, newOut, "[x] Distribute Money")
		assert.Contains(t, newOut, "[ ] Distribute Experience")
	})

	t.Run("second choice selected", func(t *testing.T) {
		m.choice = 1
		newOut := choicesView(m)
		assert.Contains(t, newOut, "[ ] Distribute Money")
		assert.Contains(t, newOut, "[x] Distribute Experience")
	})
}

func TestMoneyView(t *testing.T) {
	m := testModel(t)
	out := stripANSI(moneyView(m))

	t.Run("contains coin input fields", func(t *testing.T) {
		// textinput renders blurred fields showing only the first char of each placeholder
		assert.Contains(t, out, "Platinum")
		assert.Contains(t, out, "Gold")
		assert.Contains(t, out, "Electrum")
		assert.Contains(t, out, "Silver")
		assert.Contains(t, out, "Copper")
	})

	t.Run("shows coin priority member name", func(t *testing.T) {
		assert.Contains(t, out, "Current Coin Priority is to Keg")
	})

	t.Run("contains Submit button", func(t *testing.T) {
		assert.Contains(t, out, "Submit")
	})
}

func TestXpView(t *testing.T) {
	m := testModel(t)
	out := xpView(m)

	assert.Contains(t, out, "Xp")
	assert.Contains(t, out, "Submit")
}

func TestAddMemberView(t *testing.T) {
	m := testModel(t)
	out := stripANSI(addMemberView(m))

	t.Run("contains field placeholders", func(t *testing.T) {
		assert.Contains(t, out, "Name")
		assert.Contains(t, out, "XP")
		assert.Contains(t, out, "Platinum")
		assert.Contains(t, out, "Gold")
		assert.Contains(t, out, "Electrum")
		assert.Contains(t, out, "Silver")
		assert.Contains(t, out, "Copper")
	})

	t.Run("contains Submit button", func(t *testing.T) {
		assert.Contains(t, out, "Submit")
	})
}

func TestActivateMemberView(t *testing.T) {
	m := testModel(t)
	out := stripANSI(activateMemberView(m))

	t.Run("contains Active and Inactive headers", func(t *testing.T) {
		assert.Contains(t, out, "Active Party Members")
		assert.Contains(t, out, "Inactive Party Members")
	})

	t.Run("contains active member", func(t *testing.T) {
		assert.Contains(t, out, "Keg")
	})
}
