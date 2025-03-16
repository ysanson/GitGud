package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ysanson/GitGud/internal/git"
	"github.com/ysanson/GitGud/internal/ui"
)

type model struct {
	branches []git.Branch
	cursor   int
	logs     string
	width    int
	height   int
}

func (m model) Init() tea.Cmd {
	return ui.FetchInit
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ui.InitState:
		m.branches = msg.Branches
		m.logs = msg.Logs
	case ui.LogsMsg:
		m.logs = string(msg)
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
				return m, ui.GetGitLogs(m.branches[m.cursor])
			}
		case "down":
			if m.cursor < len(m.branches)-1 {
				m.cursor++
				return m, ui.GetGitLogs(m.branches[m.cursor])
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case ui.ErrMsg:
		fmt.Printf(msg.Error())
		return m, tea.Quit
	}
	return m, nil
}

func (m model) View() string {
	leftPanelWidth := m.width / 3
	rightPanelWidth := m.width - leftPanelWidth - 5 // Ajuster pour les bordures

	leftPanelStyle := lipgloss.NewStyle().Width(leftPanelWidth).Height(m.height - 2).Padding(1).Border(lipgloss.NormalBorder())
	rightPanelStyle := lipgloss.NewStyle().Width(rightPanelWidth).Height(m.height - 2).Padding(1).Border(lipgloss.RoundedBorder())

	var leftPanel strings.Builder
	leftPanel.WriteString("Branches:\n\n")
	for i, branch := range m.branches {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		leftPanel.WriteString(fmt.Sprintf("%s %s\n", cursor, branch.Name))
	}

	rightPanel := "Logs:\n\n" + m.logs

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPanelStyle.Render(leftPanel.String()), rightPanelStyle.Render(rightPanel))
}

func main() {
	m := model{cursor: 0}
	p := tea.NewProgram(m, tea.WithAltScreen()) // Active le mode plein écran
	if _, err := p.Run(); err != nil {
		fmt.Println("Erreur:", err)
	}
}
