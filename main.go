package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ysanson/GitGud/internal/git"
)

type model struct {
	branches []git.Branch
	cursor   int
	logs     string
	width    int
	height   int
}

type initState struct {
	branches []git.Branch
	logs     string
}

type logsMsg string

type errMsg struct{ err error }

// For messages that contain errors it's often handy to also implement the
// error interface on the message.
func (e errMsg) Error() string { return e.err.Error() }

func openGit() error {
	err := git.Open(".")
	return err
}

func getGitLogs(branch git.Branch) tea.Cmd {
	return func() tea.Msg {
		logs, err := git.ReadLogs(branch.Hash)
		if err != nil {
			return errMsg{err: err}
		}
		return logsMsg(logs)
	}
}

func fetchInit() tea.Msg {
	if err := openGit(); err != nil {
		fmt.Printf("Error: not a git repository")
		return errMsg{err: err}
	}
	branches, err := git.Branches()
	if err != nil {
		return errMsg{err: err}
	}
	logs, err := git.ReadLogs(branches[0].Hash)
	if err != nil {
		return errMsg{err: err}
	}
	return initState{branches: branches, logs: logs}
}

func (m model) Init() tea.Cmd {
	return fetchInit
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case initState:
		m.branches = msg.branches
		m.logs = msg.logs
	case logsMsg:
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
				return m, getGitLogs(m.branches[m.cursor])
			}
		case "down":
			if m.cursor < len(m.branches)-1 {
				m.cursor++
				return m, getGitLogs(m.branches[m.cursor])
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case errMsg:
		fmt.Printf(msg.Error())
		return m, tea.Quit
	}
	return m, nil
}

func (m model) View() string {
	leftPanelWidth := m.width / 3
	rightPanelWidth := m.width - leftPanelWidth - 3 // Ajuster pour les bordures

	leftPanelStyle := lipgloss.NewStyle().Width(leftPanelWidth).Height(m.height - 2).Padding(1).Border(lipgloss.NormalBorder())
	rightPanelStyle := lipgloss.NewStyle().Width(rightPanelWidth).Height(m.height - 2).Padding(1).Border(lipgloss.RoundedBorder())

	leftPanel := "Branches:\n\n"
	for i, branch := range m.branches {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		leftPanel += fmt.Sprintf("%s %s\n", cursor, branch.Name)
	}

	rightPanel := "Logs:\n\n" + m.logs

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPanelStyle.Render(leftPanel), rightPanelStyle.Render(rightPanel))
}

func main() {
	m := model{cursor: 0}
	p := tea.NewProgram(m, tea.WithAltScreen()) // Active le mode plein écran
	if _, err := p.Run(); err != nil {
		fmt.Println("Erreur:", err)
	}
}
