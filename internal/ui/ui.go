package ui

import (
	"fmt"

	"dndgoldtracker/internal/party"
	"dndgoldtracker/storage"

	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	name    = "Name"
	xp      = "XP"
	level   = "Level"
	dotChar = " • "
)

var (
	baseStyle           = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))
	subtleStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	checkboxStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	dotStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton   = focusedStyle.Render("[ Submit ]")
	blurredButton   = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	xpFields        = []string{xp}
	newMemberFields = []string{name, xp}
)

type model struct {
	activeMemberTable   table.Model
	inactiveMemberTable table.Model
	party               party.Party
	choice              int
	chosen              bool
	coinFocusIndex      int
	coinInputs          []textinput.Model
	xpFocusIndex        int
	xpInputs            []textinput.Model
	memberFocusIndex    int
	memberInputs        []textinput.Model
	virtualCursor       bool
	quitting            bool
}

// NewModel initializes the application state
func NewModel() model {
	p, err := storage.LoadParty() // Load saved data
	if err != nil {
		fmt.Println("Starting new party...")
		p = party.Party{}
	}

	newMemberFields = append(newMemberFields, party.CoinOrder...)

	amt := configureTable(p.ActiveMembers)
	imt := configureTable(p.InactiveMembers)

	ci := configureInputs(party.CoinOrder)
	xi := configureInputs(xpFields)
	mi := configureInputs(newMemberFields)

	return model{
		party:               p,
		activeMemberTable:   amt,
		inactiveMemberTable: imt,
		coinInputs:          ci,
		xpInputs:            xi,
		memberInputs:        mi,
		virtualCursor:       true,
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyPressMsg); ok {
		k := msg.String()
		if k == "ctrl+c" || (k == "q" && !m.chosen) {
			m.quitting = true
			return m, tea.Quit
		}
		if k == "esc" {
			if m.chosen {
				m.chosen = false
				m.coinFocusIndex = 0
				m.xpFocusIndex = 0
				m.memberFocusIndex = 0
				resetInputs(m.coinInputs)
				resetInputs(m.xpInputs)
				resetInputs(m.memberInputs)
				blurTable(&m.activeMemberTable)
				return m, nil
			}
			m.quitting = true
			return m, tea.Quit
		}
	}

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	if !m.chosen {
		return updateChoices(msg, m)
	}

	switch m.choice {
	case 0:
		return updateMoney(msg, m)
	case 1:
		return updateExperience(msg, m)
	case 2:
		return updateAddMember(msg, m)
	case 3:
		return updateActivateMembers(msg, m)
	default:
		return m, nil
	}
}

// The main view, which just calls the appropriate sub-view
func (m model) View() tea.View {
	var s string
	if m.quitting {
		return tea.NewView("\n  See you later!\n\n")
	}

	if !m.chosen {
		s = choicesView(m)
	} else {
		switch m.choice {
		case 0:
			s = moneyView(m)
		case 1:
			s = xpView(m)
		case 2:
			s = addMemberView(m)
		case 3:
			s = activateMemberView(m)
		default:
			s = "Don't do that"
		}
	}

	return tea.NewView(baseStyle.Render("\n" + s + "\n\n"))
}
