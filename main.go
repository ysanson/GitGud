package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ysanson/GitGud/internal/git"
	"github.com/ysanson/GitGud/internal/ui"
)

type fileStatus map[string]string

type model struct {
	branches   []git.Branch
	untracked  []string
	modified   fileStatus
	staged     fileStatus
	cursor     int
	logs       string
	width      int
	height     int
	currentTab int
}

func (m model) Init() tea.Cmd {
	return ui.GetGitStatus()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ui.LogsState:
		m.branches = msg.Branches
		m.logs = msg.Logs
	case ui.LogsMsg:
		m.logs = string(msg)
	case ui.StatusState:
		m.untracked = msg.Untracked
		m.modified = msg.Modified
		m.staged = msg.Staged
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		switch msg.String() {
		case "tab":
			m.currentTab = (m.currentTab + 1) % 2
			if m.currentTab == 0 {
				return m, ui.GetGitStatus()
			} else if m.currentTab == 1 {
				return m, ui.FetchInit
			}
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

func (m model) statusView() string {
	style := lipgloss.NewStyle().Underline(true)
	untrTitle := style.Render("Untracked:\n\n")
	modifTitle := style.Render("Modified:\n\n")
	stagedTitle := style.Foreground(lipgloss.Color("#ACD8A9")).Render("Staged:\n\n")
	var untracked strings.Builder
	if len(m.untracked) == 0 {
		untracked.WriteString("(no files)")
	} else {
		for _, f := range m.untracked {
			untracked.WriteString("- " + f + "\n")
		}
	}
	status := func(files map[string]string) string {
		var sb strings.Builder
		if len(files) == 0 {
			sb.WriteString("(no files)")
		} else {
			for path, status := range files {
				sb.WriteString("- " + path + ": " + status + "\n")
			}
		}
		return sb.String()
	}
	thirdWidth := m.width/3 - 3
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Width(thirdWidth).Padding(1).Border(lipgloss.NormalBorder()).Render(untrTitle+untracked.String()),
		lipgloss.NewStyle().Width(thirdWidth).Padding(1).Border(lipgloss.NormalBorder()).Render(modifTitle+status(m.modified)),
		lipgloss.NewStyle().Width(thirdWidth).Padding(1).Border(lipgloss.NormalBorder()).Render(stagedTitle+status(m.staged)),
	)
}

func (m model) logsView() string {
	leftPanelWidth := m.width / 3
	rightPanelWidth := m.width - leftPanelWidth - 5

	leftPanelStyle := lipgloss.NewStyle().Width(leftPanelWidth).Height(m.height - 10).Padding(1).Border(lipgloss.NormalBorder())
	rightPanelStyle := lipgloss.NewStyle().Width(rightPanelWidth).Height(m.height - 10).Padding(1).Border(lipgloss.NormalBorder())

	leftPanel := "Branches:\n\n"
	for i, branch := range m.branches {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		leftPanel += fmt.Sprintf("%s %s\n", cursor, branch.Name)
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanelStyle.Render(leftPanel),
		rightPanelStyle.Render("Logs:\n\n"+m.logs),
	)
}

func (m model) View() string {
	tabs := []string{"Status", "Logs"}
	tabHeader := ""
	for i, t := range tabs {
		if i == m.currentTab {
			tabHeader += lipgloss.NewStyle().Bold(true).Underline(true).Render(t) + "   "
		} else {
			tabHeader += t + "   "
		}
	}

	var content string
	switch m.currentTab {
	case 0:
		content = m.statusView()
	case 1:
		content = m.logsView()
	}

	return tabHeader + "\n\n" + content
}

func main() {
	m := model{cursor: 0, currentTab: 0}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}
