package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ysanson/GitGud/internal/git"
)

type InitState struct {
	Branches []git.Branch
	Logs     string
}

type LogsMsg string

type ErrMsg struct{ err error }

// For messages that contain errors it's often handy to also implement the
// error interface on the message.
func (e ErrMsg) Error() string { return e.err.Error() }

func openGit() error {
	err := git.Open(".")
	return err
}

func GetGitLogs(branch git.Branch) tea.Cmd {
	return func() tea.Msg {
		logs, err := git.ReadLogs(branch.Hash)
		if err != nil {
			return ErrMsg{err: err}
		}
		return LogsMsg(logs)
	}
}

func FetchInit() tea.Msg {
	if err := openGit(); err != nil {
		fmt.Printf("Error: not a git repository")
		return ErrMsg{err: err}
	}
	branches, err := git.Branches()
	if err != nil {
		return ErrMsg{err: err}
	}
	logs, err := git.ReadLogs(branches[0].Hash)
	if err != nil {
		return ErrMsg{err: err}
	}
	return InitState{Branches: branches, Logs: logs}
}
