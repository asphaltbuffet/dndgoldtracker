package ui

import (
	"log/slog"
	"slices"
	"strconv"

	"dndgoldtracker/internal/party"
	"dndgoldtracker/internal/storage"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
)

// Sub-Updates

// Update loop for the first view where you're choosing a task.
func updateChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "j", "down":
			m.choice++
			if m.choice > 3 {
				m.choice = 0
			}
		case "k", "up":
			m.choice--
			if m.choice < 0 {
				m.choice = 3
			}
		case "enter":
			m.chosen = true
			if m.choice == 3 {
				focusTable(&m.activeMemberTable)
				blurTable(&m.inactiveMemberTable)
				m.activeMemberTable.SetCursor(0)
			}
			return m, nil
		}
	}

	return m, nil
}

// Update loop for updating party money
func updateMoney(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		// Change cursor mode
		case "ctrl+r":
			changeCursorMode(m.coinInputs, &m.virtualCursor)
			return m, nil

		// Set focus to next input
		case "enter":
			// Did the user press enter while the submit button was focused?
			// If so, Distribute money.
			if m.coinFocusIndex == len(m.coinInputs) {
				var err error
				coinMap := make(map[party.Coin]int)
				// Set any unset values to 0
				handleUnsetInputs(m.coinInputs)

				for i := range party.CoinOrder {
					coinMap[party.CoinOrder[i]], err = strconv.Atoi(m.coinInputs[i].Value())
					if err != nil {
						slog.Error("invalid coin value", "type", party.CoinOrder[i], "input", m.coinInputs[i].Value(), "err", err)

						return m, nil
					}

					slog.Debug("coin input parsed", "type", party.CoinOrder[i], "amount", coinMap[party.CoinOrder[i]])
				}

				// Distribute the coins to the party
				party.DistributeCoins(&m.party, coinMap)
				saveUpdateReset(&m)

				m.chosen = false
				return m, nil
			}
			// Cycle indexes
		case "tab", "down":
			m.coinFocusIndex++
			slog.Debug("focus changed", "index", m.coinFocusIndex, "input", "coin")
			cmds := updateFocusIndex(&m.coinFocusIndex, m.coinInputs)
			return m, tea.Batch(cmds...)
		case "up", "shift+tab":
			m.coinFocusIndex--
			slog.Debug("focus changed", "index", m.coinFocusIndex, "input", "coin")
			cmds := updateFocusIndex(&m.coinFocusIndex, m.coinInputs)
			return m, tea.Batch(cmds...)
		}
	}
	// Handle character input and blinking
	var cmd tea.Cmd
	m.coinInputs, cmd = updateInputs(msg, m.coinInputs)

	return m, cmd
}

// Update loop for updating party experience
func updateExperience(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		// Change cursor mode
		case "ctrl+r":
			changeCursorMode(m.xpInputs, &m.virtualCursor)
			return m, nil

		// Set focus to next input
		case "enter":
			// Did the user press enter while the submit button was focused?
			// If so, Distribute xp.
			if m.xpFocusIndex == len(m.xpInputs) {
				var err error
				handleUnsetInputs(m.xpInputs)

				xp, err := strconv.Atoi(m.xpInputs[0].Value())
				if err != nil {
					slog.Error("invalid XP input", "input", m.xpInputs[0].Value(), "err", err)

					return m, nil
				}

				slog.Debug("parsed xp input", "xp", xp)

				m.party.DistributeExperience(xp)
				saveUpdateReset(&m)

				m.chosen = false
				return m, nil
			}

		case "tab", "down":
			m.xpFocusIndex++
			cmds := updateFocusIndex(&m.xpFocusIndex, m.xpInputs)
			return m, tea.Batch(cmds...)

		case "up", "shift+tab":
			m.xpFocusIndex--
			cmds := updateFocusIndex(&m.xpFocusIndex, m.xpInputs)
			return m, tea.Batch(cmds...)
		}
	}
	// Handle character input and blinking
	var cmd tea.Cmd
	m.xpInputs, cmd = updateInputs(msg, m.xpInputs)

	return m, cmd
}

// Update loop for adding members
func updateAddMember(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		// Change cursor mode
		case "ctrl+r":
			changeCursorMode(m.memberInputs, &m.virtualCursor)
			return m, nil

		// Set focus to next input
		case "enter":
			// Did the user press enter while the submit button was focused?
			// If so, Distribute money.
			if m.memberFocusIndex == len(m.memberInputs) {
				var err error
				if m.memberInputs[0].Value() == "" {
					slog.Warn("member name is required")
					m.chosen = false
					return m, nil
				}
				name := m.memberInputs[0].Value()
				// Set any unset values other than name to 0
				handleUnsetInputs(m.memberInputs)
				xp, err := strconv.Atoi(m.memberInputs[1].Value())
				if err != nil {
					slog.Warn("invalid XP value for new member")
					return m, nil
				}

				newMemberCoins := m.memberInputs[2:len(m.memberInputs)]
				newMemberMoney := make(map[party.Coin]int)
				for i := range newMemberCoins {
					newMemberMoney[party.Coin(i)], err = strconv.Atoi(newMemberCoins[i].Value())
					if err != nil {
						slog.Warn("invalid coin value for new member", "field", newMemberCoins[i].Placeholder)
						return m, nil
					}
				}

				// Add the new party Member
				m.party.AddMember(name, xp, newMemberMoney)
				saveUpdateReset(&m)

				m.chosen = false
				return m, nil
			}
		// Cycle indexes
		case "tab", "down":
			m.memberFocusIndex++
			slog.Debug("focus changed", "index", m.memberFocusIndex, "input", "member")
			cmds := updateFocusIndex(&m.memberFocusIndex, m.memberInputs)
			return m, tea.Batch(cmds...)
		case "up", "shift+tab":
			m.memberFocusIndex--
			slog.Debug("focus changed", "index", m.memberFocusIndex, "input", "member")
			cmds := updateFocusIndex(&m.memberFocusIndex, m.memberInputs)
			return m, tea.Batch(cmds...)
		}
	}
	// Handle character input and blinking
	var cmd tea.Cmd
	m.memberInputs, cmd = updateInputs(msg, m.memberInputs)

	return m, cmd
}

// Update loop for activating or deactivating members
func updateActivateMembers(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var activeCmd tea.Cmd
	var inactiveCmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "tab":
			// Change table focus with tab
			if m.activeMemberTable.Focused() {
				blurTable(&m.activeMemberTable)
				focusTable(&m.inactiveMemberTable)
				m.inactiveMemberTable.SetCursor(0)
			} else {
				focusTable(&m.activeMemberTable)
				m.activeMemberTable.SetCursor(0)
				blurTable(&m.inactiveMemberTable)
			}
		case "enter":
			var selectedTable *table.Model
			var selectedMembers *[]party.Member
			var unselectedMembers *[]party.Member
			// Move the selected member from their current table to the new one
			if m.activeMemberTable.Focused() {
				selectedTable = &m.activeMemberTable
				selectedMembers = &m.party.ActiveMembers
				unselectedMembers = &m.party.InactiveMembers
			} else {
				selectedTable = &m.inactiveMemberTable
				selectedMembers = &m.party.InactiveMembers
				unselectedMembers = &m.party.ActiveMembers
			}
			// activate/deactivate member
			if len(selectedTable.SelectedRow()) == 0 {
				// Unselected cursor or empty table
				// Set cursor to first element and return
				slog.Debug("no row selected, resetting cursor")
				selectedTable.SetCursor(0)
				return m, nil
			}

			memberName := selectedTable.SelectedRow()[0]
			from, to := "Active", "Inactive"
			if !m.activeMemberTable.Focused() {
				from, to = "Inactive", "Active"
			}
			slog.Info("moving member", "name", memberName, "from", from, "to", to)

			memberIndex := slices.IndexFunc(*selectedMembers, func(m party.Member) bool { return m.Name == memberName })
			party.ChangeMemberGroup(selectedMembers, unselectedMembers, memberIndex)
			m.activeMemberTable.SetRows(membersToRows(m.party.ActiveMembers))
			m.inactiveMemberTable.SetRows(membersToRows(m.party.InactiveMembers))
		case "s":
			err := storage.SaveParty(&m.party)
			if err != nil {
				slog.Error("failed to save party data", "err", err)
			}

			blurTable(&m.activeMemberTable)
			m.chosen = false
		}
	}

	m.activeMemberTable, activeCmd = m.activeMemberTable.Update(msg)
	m.inactiveMemberTable, inactiveCmd = m.inactiveMemberTable.Update(msg)
	return m, tea.Batch(activeCmd, inactiveCmd)
}
