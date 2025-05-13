package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list     list.Model
	selected string
	viewport viewport.Model
	showMenu bool
}

func initialModel() model {
	items := []list.Item{
		item{title: "Personal data", desc: "Name, who to call if..."},
		item{title: "Check methods", desc: "medical, email, open check, scrapper, others"},
		item{title: "Where to run this", desc: "only here, self-hosted web version"},
		item{title: "Just kill me", desc: "closes the CLI"},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "stillAlive - Long-term Life Monitor"

	vp := viewport.New(0, 0)
	vp.Style = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder())

	return model{
		list:     l,
		viewport: vp,
		showMenu: true,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.selected = i.title
				if m.selected == "Just kill me" {
					return m, tea.Quit
				}
				m.showMenu = false
			}
			return m, nil
		case "esc":
			if !m.showMenu {
				m.showMenu = true
				return m, nil
			}
		}
	case tea.WindowSizeMsg:
		h, v := lipgloss.NewStyle().Margin(1, 2).GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.viewport.Width = msg.Width - h
		m.viewport.Height = msg.Height - v
	}

	if m.showMenu {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.showMenu {
		return m.list.View()
	}

	content := ""
	switch m.selected {
	case "Personal data":
		content = "Personal Data Configuration\n\n" +
			"Full name: _____________\n" +
			"Emergency contact: _____________\n" +
			"Phone: _____________\n" +
			"Relationship: _____________\n\n" +
			"Press ESC to return"
	case "Check methods":
		content = "Verification Methods\n\n" +
			"[ ] Medical check\n" +
			"[ ] Email verification\n" +
			"[ ] Web scrapper\n" +
			"[ ] Others\n\n" +
			"Press ESC to return"
	case "Where to run this":
		content = "Execution Options\n\n" +
			"[ ] Local only\n" +
			"[ ] Self-hosted web version\n\n" +
			"Press ESC to return"
	}

	m.viewport.SetContent(content)
	return m.viewport.View()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error starting stillAlive: %v", err)
		os.Exit(1)
	}
}