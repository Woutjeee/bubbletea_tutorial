package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
)

type focusState int

const (
	focusList focusState = iota
	focusInput
)

type model struct {
	choises    []string
	cursor     int
	selected   map[int]struct{}
	textInput  textinput.Model
	err        error
	focusIndex focusState
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Grocery.."
	ti.Focus()
	ti.CharLimit = 150
	ti.Width = 20

	return model{
		choises:    []string{},
		selected:   make(map[int]struct{}),
		textInput:  ti,
		err:        nil,
		focusIndex: focusInput,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.focusIndex == focusList {
				m.focusIndex = focusInput
				m.textInput.Focus()
			} else {
				m.focusIndex = focusList
				m.textInput.Blur()
			}
		case "up", "k":
			if m.focusIndex == focusList && m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.focusIndex == focusList && m.cursor < len(m.choises)-1 {
				m.cursor++
			}
		case "enter", " ":
			if m.focusIndex == focusList {
				_, ok := m.selected[m.cursor]
				if ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}
				}
			} else if m.focusIndex == focusInput {
				m.choises = append(m.choises, m.textInput.Value())
				m.textInput.Reset()
			}
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	if m.focusIndex == focusInput {
		m.textInput, cmd = m.textInput.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	s := "What shoud we buy at the market?\n\n"

	s += fmt.Sprintf("What more should we get?\n%s\n\n", m.textInput.View())

	if len(m.choises) != 0 {
		for i, choice := range m.choises {
			cursor := " "
			if m.focusIndex == focusList && m.cursor == i {
				cursor = ">"
			}

			checked := " "
			if _, ok := m.selected[i]; ok {
				checked = "x"
			}
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
		}
	}

	if m.focusIndex == focusList {
		s += "\n\nUse ↑/↓ to navigate, Space to select, Tab to focus on the input, Enter to add items, and q to quit.\n"
	} else {
		s += "\n\nType your item, Tab to focus on the list, and q to quit.\n"
	}
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
