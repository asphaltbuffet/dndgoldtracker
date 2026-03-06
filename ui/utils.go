package ui

import (
	"fmt"
	"strconv"
	"strings"

	"dndgoldtracker/models"
	"dndgoldtracker/storage"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func checkbox(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}

// Convert members to table rows
func membersToRows(members []models.Member) []table.Row {
	var rows []table.Row
	for _, m := range members {
		rows = append(rows, table.Row{
			m.Name,
			strconv.Itoa(m.XP),
			strconv.Itoa(m.Level),
			strconv.Itoa(m.Coins[models.Platinum]),
			strconv.Itoa(m.Coins[models.Gold]),
			strconv.Itoa(m.Coins[models.Electrum]),
			strconv.Itoa(m.Coins[models.Silver]),
			strconv.Itoa(m.Coins[models.Copper]),
		})
	}
	return rows
}

func (m *model) updateInputs(msg tea.Msg, inputs []textinput.Model) tea.Cmd {
	cmds := make([]tea.Cmd, len(inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range inputs {
		inputs[i], cmds[i] = inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func configureInputs(placeholders []string) []textinput.Model {
	i := make([]textinput.Model, len(placeholders))
	var t textinput.Model
	for j := range i {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32
		t.Placeholder = placeholders[j]

		// Focus the first element
		if j == 0 {
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		}

		i[j] = t
	}

	return i
}

func configureTable(members []models.Member) table.Model {
	columns := []table.Column{
		{Title: name, Width: 10},
		{Title: xp, Width: 6},
		{Title: level, Width: 6},
		{Title: models.Platinum, Width: 10},
		{Title: models.Gold, Width: 6},
		{Title: models.Electrum, Width: 10},
		{Title: models.Silver, Width: 8},
		{Title: models.Copper, Width: 8},
	}

	rows := membersToRows(members)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(5),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t
}

func updateTableData(members []models.Member, t *table.Model) *table.Model {
	rows := membersToRows(members)
	t.SetRows(rows)
	return t
}

func resetInputs(inputs []textinput.Model) {
	for i := range inputs {
		inputs[i].Reset()
	}
}

func buildInputList(inputs []textinput.Model, focusIndex int, cursorMode cursor.Mode) string {
	var msg strings.Builder
	for i := range inputs {
		msg.WriteString(inputs[i].View())
		if i < len(inputs)-1 {
			msg.WriteRune('\n')
		}
	}

	button := &blurredButton
	if focusIndex == len(inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&msg, "\n\n%s\n\n", *button)

	msg.WriteString(helpStyle.Render("cursor mode is "))
	msg.WriteString(cursorModeHelpStyle.Render(cursorMode.String()))
	msg.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return msg.String()
}

func handleUnsetInputs(inputs []textinput.Model) {
	for i := range inputs {
		if inputs[i].Value() == "" {
			inputs[i].SetValue("0")
		}
	}
}

func saveUpdateReset(m *model) {
	_ = storage.SaveParty(&m.party)
	updateTableData(m.party.ActiveMembers, &m.activeMemberTable)
	resetInputs(m.coinInputs)
	resetInputs(m.xpInputs)
	resetInputs(m.memberInputs)
}

func changeCursorMode(inputs []textinput.Model, cursorMode *cursor.Mode) []tea.Cmd {
	*cursorMode++
	if *cursorMode > cursor.CursorHide {
		*cursorMode = cursor.CursorBlink
	}
	cmds := make([]tea.Cmd, len(inputs))
	for i := range inputs {
		cmds[i] = inputs[i].Cursor.SetMode(*cursorMode)
	}
	return cmds
}

func updateFocusIndex(focusIndex *int, inputs []textinput.Model) []tea.Cmd {
	if *focusIndex > len(inputs) {
		*focusIndex = 0
	} else if *focusIndex < 0 {
		*focusIndex = len(inputs)
	}

	cmds := make([]tea.Cmd, len(inputs))
	for i := 0; i <= len(inputs)-1; i++ {
		if i == *focusIndex {
			// Set focused state
			cmds[i] = inputs[i].Focus()
			inputs[i].PromptStyle = focusedStyle
			inputs[i].TextStyle = focusedStyle
			continue
		}
		// Remove focused state
		inputs[i].Blur()
		inputs[i].PromptStyle = noStyle
		inputs[i].TextStyle = noStyle
	}

	return cmds
}
