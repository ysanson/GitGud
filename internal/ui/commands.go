package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ysanson/GitGud/internal/git"
)

type LogsState struct {
	Branches []git.Branch
	Logs     string
}

type StatusState struct {
	Untracked []string
	Modified  map[string]string
	Staged    map[string]string
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
	return LogsState{Branches: branches, Logs: logs}
}

func GetGitStatus() tea.Cmd {
	return func() tea.Msg {
		if err := openGit(); err != nil {
			fmt.Printf("Error: not a git repository")
			return ErrMsg{err: err}
		}
		untracked, modified, staged, err := git.FetchGitStatus()
		if err != nil {
			return ErrMsg{err: err}
		}
		return StatusState{Untracked: untracked, Modified: modified, Staged: staged}
	}
}
