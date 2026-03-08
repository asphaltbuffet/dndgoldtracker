package ui

import (
	"dndgoldtracker/commands"
	"fmt"
	"strings"
)

// Sub-Views

// The first view, where you're choosing a task
func choicesView(m model) string {
	choice := m.choice
	var msg string
	if (len(m.party.ActiveMembers)) > 0 {
		msg += "Current Coin Priority is to " +
			focusedStyle.Render(m.party.ActiveMembers[commands.GetFirstCoinPriority(&m.party)].Name) + "\n"
	}

	msg += baseStyle.Render(m.activeMemberTable.View())

	msg += "\nWhat would you like to do?"
	msg += "\n"

	msg += fmt.Sprintf(
		"\n%s\n%s\n%s\n%s\n",
		checkbox("Distribute Money", choice == 0),
		checkbox("Distribute Experience", choice == 1),
		checkbox("Add Member", choice == 2),
		checkbox("Activate/Deactivate Party Members", choice == 3),
	)

	msg += subtleStyle.Render("j/k, up/down: select") + dotStyle +
		subtleStyle.Render("enter: choose") + dotStyle +
		subtleStyle.Render("q, esc: quit")

	return msg
}

// The view for adding money
func moneyView(m model) string {
	var msg strings.Builder

	msg.WriteString(baseStyle.Render("Money entered here will be distributed to all party members as equally as possible.\n" +
		"Extra coins are distributed based on a priority system that rotates.\n"))
	if (len(m.party.ActiveMembers)) > 0 {
		msg.WriteString("\nCurrent Coin Priority is to " +
			focusedStyle.Render(m.party.ActiveMembers[commands.GetFirstCoinPriority(&m.party)].Name))
	}

	msg.WriteString("\n" + buildInputList(m.coinInputs, m.coinFocusIndex, m.virtualCursor))
	return msg.String()
}

// The view for adding experience
func xpView(m model) string {
	var msg strings.Builder
	msg.WriteString("Xp entered here will be distributed to all party members equally\n\n")
	msg.WriteString(buildInputList(m.xpInputs, m.xpFocusIndex, m.virtualCursor))
	return msg.String()
}

func addMemberView(m model) string {
	var msg strings.Builder
	msg.WriteString("Enter the new party member's data\n")
	msg.WriteString(buildInputList(m.memberInputs, m.memberFocusIndex, m.virtualCursor))
	return msg.String()
}

func activateMemberView(m model) string {
	var msg strings.Builder
	if m.activeMemberTable.Focused() {
		msg.WriteString(focusedStyle.Render("Active Party Members\n"))
		msg.WriteString(focusedStyle.Render("\n" + m.activeMemberTable.View()))
		msg.WriteString(blurredStyle.Render("\nInactive Party Members\n"))
		msg.WriteString(blurredStyle.Render("\n" + m.inactiveMemberTable.View()))
	} else {
		msg.WriteString(blurredStyle.Render("Active Party Members\n"))
		msg.WriteString(blurredStyle.Render("\n" + m.activeMemberTable.View()))
		msg.WriteString(focusedStyle.Render("\nInactive Party Members\n"))
		msg.WriteString(focusedStyle.Render("\n" + m.inactiveMemberTable.View()))
	}

	msg.WriteString(subtleStyle.Render("\nup/down: select") + dotStyle +
		subtleStyle.Render("enter: activate/deactivate member") + dotStyle +
		subtleStyle.Render("s: save and return to menu") + dotStyle +
		subtleStyle.Render("tab: switch table"))
	return msg.String()
}
