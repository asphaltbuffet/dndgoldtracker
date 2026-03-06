package ui

import (
	"fmt"
	"log"
	"slices"
	"strconv"

	"dndgoldtracker/commands"
	"dndgoldtracker/models"
	"dndgoldtracker/storage"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

// Sub-Updates

// Update loop for the first view where you're choosing a task.
func updateChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
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
	case tea.KeyMsg:
		switch msg.String() {
		// Change cursor mode
		case "ctrl+r":
			var cmds []tea.Cmd
			cmds = changeCursorMode(m.xpInputs, &m.cursorMode)

			return m, tea.Batch(cmds...)
		// Set focus to next input
		case "enter":
			// Did the user press enter while the submit button was focused?
			// If so, Distribute money.
			if m.coinFocusIndex == len(m.coinInputs) {
				var err error
				coinMap := make(map[string]int)
				// Set any unset values to 0
				handleUnsetInputs(m.coinInputs)

				for i := range models.CoinOrder {
					coinMap[models.CoinOrder[i]], err = strconv.Atoi(m.coinInputs[i].Value())
					log.Printf("CoinMap entry for %s: %d\n", models.CoinOrder[i], coinMap[models.CoinOrder[i]])
					if err != nil {
						fmt.Printf("Invalid input for %s, try again\n", models.CoinOrder[i])
						return m, nil
					}
				}

				// Distribute the coins to the party
				commands.DistributeCoins(&m.party, coinMap)
				saveUpdateReset(&m)

				m.chosen = false
				return m, nil
			}
			// Cycle indexes
		case "tab", "down":
			m.coinFocusIndex++
			log.Printf("coinFocusIndex = %d", m.coinFocusIndex)
			cmds := updateFocusIndex(&m.coinFocusIndex, m.coinInputs)
			return m, tea.Batch(cmds...)
		case "up", "shift+tab":
			m.coinFocusIndex--
			log.Printf("coinFocusIndex = %d", m.coinFocusIndex)
			cmds := updateFocusIndex(&m.coinFocusIndex, m.coinInputs)
			return m, tea.Batch(cmds...)
		}
	}
	// Handle character input and blinking
	cmd := m.updateInputs(msg, m.coinInputs)

	return m, cmd
}

// Update loop for updating party experience
func updateExperience(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Change cursor mode
		case "ctrl+r":
			var cmds []tea.Cmd
			cmds = changeCursorMode(m.xpInputs, &m.cursorMode)
			return m, tea.Batch(cmds...)

		// Set focus to next input
		case "enter":
			// Did the user press enter while the submit button was focused?
			// If so, Distribute xp.
			if m.xpFocusIndex == len(m.xpInputs) {
				var err error
				handleUnsetInputs(m.xpInputs)

				xp, err := strconv.Atoi(m.xpInputs[0].Value())
				if err != nil {
					log.Println("Invalid input for experience, try again")
					return m, nil
				}

				commands.DistributeExperience(&m.party, xp)
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
	cmd := m.updateInputs(msg, m.xpInputs)

	return m, cmd
}

// Update loop for adding members
func updateAddMember(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Change cursor mode
		case "ctrl+r":
			var cmds []tea.Cmd
			cmds = changeCursorMode(m.memberInputs, &m.cursorMode)

			return m, tea.Batch(cmds...)
		// Set focus to next input
		case "enter":
			// Did the user press enter while the submit button was focused?
			// If so, Distribute money.
			if m.memberFocusIndex == len(m.memberInputs) {
				var err error
				if m.memberInputs[0].Value() == "" {
					log.Println("Name value required, try again")
					m.chosen = false
					return m, nil
				}
				name := m.memberInputs[0].Value()
				// Set any unset values other than name to 0
				handleUnsetInputs(m.memberInputs)
				xp, err := strconv.Atoi(m.memberInputs[1].Value())
				if err != nil {
					log.Println("Error occurred with xp, try again")
					return m, nil
				}

				newMemberCoins := m.memberInputs[2:len(m.memberInputs)]
				newMemberMoney := make(map[string]int)
				for i := range newMemberCoins {
					newMemberMoney[newMemberCoins[i].Placeholder], err = strconv.Atoi(newMemberCoins[i].Value())
					if err != nil {
						log.Printf("Error occurred with %s. Please try again", newMemberCoins[i].Placeholder)
						return m, nil
					}
				}

				// Add the new party Member
				commands.AddMember(&m.party, name, xp, newMemberMoney)
				saveUpdateReset(&m)

				m.chosen = false
				return m, nil
			}
		// Cycle indexes
		case "tab", "down":
			m.memberFocusIndex++
			log.Printf("memberFocusIndex = %d", m.memberFocusIndex)
			cmds := updateFocusIndex(&m.memberFocusIndex, m.memberInputs)
			return m, tea.Batch(cmds...)
		case "up", "shift+tab":
			m.memberFocusIndex--
			log.Printf("memberFocusIndex = %d", m.memberFocusIndex)
			cmds := updateFocusIndex(&m.memberFocusIndex, m.memberInputs)
			return m, tea.Batch(cmds...)
		}
	}
	// Handle character input and blinking
	cmd := m.updateInputs(msg, m.memberInputs)

	return m, cmd
}

// Update loop for activating or deactivating members
func updateActivateMembers(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var activeCmd tea.Cmd
	var inactiveCmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
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
			var selectedMembers *[]models.Member
			var unselectedMembers *[]models.Member
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
			if len(selectedTable.SelectedRow()) <= 0 {
				// Unselected cursor or empty table
				// Set cursor to first element and return
				log.Println("Unselected cursor thing")
				selectedTable.SetCursor(0)
				return m, nil
			}

			memberName := selectedTable.SelectedRow()[0]
			if m.activeMemberTable.Focused() {
				log.Printf("Moving %s from %s to %s", memberName, "Active", "Inactive")
			} else {
				log.Printf("Moving %s from %s to %s", memberName, "Inactive", "Active")
			}

			memberIndex := slices.IndexFunc(*selectedMembers, func(m models.Member) bool { return m.Name == memberName })
			commands.ChangeMemberGroup(selectedMembers, unselectedMembers, memberIndex)
			m.activeMemberTable.SetRows(membersToRows(m.party.ActiveMembers))
			m.inactiveMemberTable.SetRows(membersToRows(m.party.InactiveMembers))
		case "s":
			_ = storage.SaveParty(&m.party)
			blurTable(&m.activeMemberTable)
			m.chosen = false
		}
	}

	m.activeMemberTable, activeCmd = m.activeMemberTable.Update(msg)
	m.inactiveMemberTable, inactiveCmd = m.inactiveMemberTable.Update(msg)
	return m, tea.Batch(activeCmd, inactiveCmd)
}
