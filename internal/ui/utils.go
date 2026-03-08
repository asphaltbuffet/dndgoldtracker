package ui

import (
	"fmt"
	"strconv"
	"strings"

	"dndgoldtracker/internal/party"
	"dndgoldtracker/internal/storage"

	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func checkbox(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}

// Convert members to table rows
func membersToRows(members []party.Member) []table.Row {
	rows := make([]table.Row, 0, len(members))
	for _, m := range members {
		rows = append(rows, table.Row{
			m.Name,
			strconv.Itoa(m.XP),
			strconv.Itoa(m.Level),
			strconv.Itoa(m.Coins[party.Platinum]),
			strconv.Itoa(m.Coins[party.Gold]),
			strconv.Itoa(m.Coins[party.Electrum]),
			strconv.Itoa(m.Coins[party.Silver]),
			strconv.Itoa(m.Coins[party.Copper]),
		})
	}
	return rows
}

func (m model) updateInputs(msg tea.Msg, inputs []textinput.Model) tea.Cmd {
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
		t.CharLimit = 32
		t.Placeholder = placeholders[j]
		t.SetWidth(20)

		s := textinput.DefaultStyles(false)

		// Focus the first element
		if j == 0 {
			t.Focus()
			s.Focused.Prompt = focusedStyle
			s.Focused.Text = focusedStyle
		}

		t.SetStyles(s)
		i[j] = t
	}

	return i
}

func configureTable(members []party.Member) table.Model {
	columns := []table.Column{
		{Title: name, Width: 10},
		{Title: xp, Width: 6},
		{Title: level, Width: 6},
		{Title: party.Platinum.String(), Width: 10},
		{Title: party.Gold.String(), Width: 6},
		{Title: party.Electrum.String(), Width: 10},
		{Title: party.Silver.String(), Width: 8},
		{Title: party.Copper.String(), Width: 8},
	}

	rows := membersToRows(members)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(5),
		table.WithWidth(80),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = lipgloss.NewStyle()
	t.SetStyles(s)

	return t
}

func blurTable(t *table.Model) {
	t.Blur()
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = lipgloss.NewStyle()
	t.SetStyles(s)
}

func focusTable(t *table.Model) {
	t.Focus()
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
}

func updateTableData(members []party.Member, t *table.Model) *table.Model {
	rows := membersToRows(members)
	t.SetRows(rows)
	return t
}

func resetInputs(inputs []textinput.Model) {
	for i := range inputs {
		inputs[i].Reset()
	}
}

func buildInputList(inputs []textinput.Model, focusIndex int, virtualCursor bool) string {
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

	cursorModeStr := "real"
	if virtualCursor {
		cursorModeStr = "virtual"
	}
	msg.WriteString(helpStyle.Render("cursor mode is "))
	msg.WriteString(cursorModeHelpStyle.Render(cursorModeStr))
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

// changeCursorMode toggles between virtual (blinking) and real terminal cursor.
func changeCursorMode(inputs []textinput.Model, virtualCursor *bool) {
	*virtualCursor = !*virtualCursor
	for i := range inputs {
		inputs[i].SetVirtualCursor(*virtualCursor)
	}
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
			s := inputs[i].Styles()
			s.Focused.Prompt = focusedStyle
			s.Focused.Text = focusedStyle
			inputs[i].SetStyles(s)
			continue
		}
		// Remove focused state
		inputs[i].Blur()
		s := inputs[i].Styles()
		s.Focused.Prompt = noStyle
		s.Focused.Text = noStyle
		inputs[i].SetStyles(s)
	}

	return cmds
}
