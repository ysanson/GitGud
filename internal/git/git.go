package git

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Branch struct {
	Name string
	Hash plumbing.Hash
}

var (
	repo *git.Repository
)

func Open(path string) error {
	r, err := git.PlainOpen(path)
	if err != nil {
		fmt.Println(err)
	}
	repo = r
	return err
}

func Logs(hash plumbing.Hash) (object.CommitIter, error) {
	if repo != nil {
		return repo.Log(&git.LogOptions{From: hash})
	}
	return nil, git.ErrRepositoryNotExists
}

func Branches() ([]Branch, error) {
	if repo != nil {
		branches, err := repo.Branches()
		if err != nil {
			return nil, err
		}
		br := make([]Branch, 0, 5)
		branches.ForEach(func(r *plumbing.Reference) error {
			br = append(br, Branch{Name: r.Name().Short(), Hash: r.Hash()})
			return nil
		})
		return br, nil
	}
	return nil, git.ErrRepositoryNotExists
}

func ReadLogs(hash plumbing.Hash) (string, error) {
	logs, err := Logs(hash)
	if err != nil {
		return "", err
	}
	var str strings.Builder
	logs.ForEach(func(c *object.Commit) error {
		str.WriteString(fmt.Sprintf("%s: %s\n", c.Hash.String()[:7], c.Message))
		return nil
	})
	return str.String(), nil
}

func getStatus(status git.StatusCode) string {
	switch status {
	case git.Unmodified:
		return "Unmodified"
	case git.Untracked:
		return "Untracked"
	case git.Modified:
		return "Modified"
	case git.Added:
		return "Added"
	case git.Deleted:
		return "Deleted"
	case git.Renamed:
		return "Renamed"
	case git.Copied:
		return "Copied"
	case git.UpdatedButUnmerged:
		return "UpdatedButUnmerged"
	default:
		return "Unknown"
	}
}

func FetchGitStatus() (untracked []string, modified, staged map[string]string, e error) {
	if repo == nil {
		return
	}
	w, err := repo.Worktree()
	if err != nil {
		return untracked, modified, staged, err
	}

	status, err := w.Status()
	if err != nil {
		return untracked, modified, staged, err
	}
	if status.IsClean() {
		return
	}

	modified = make(map[string]string)
	staged = make(map[string]string)

	for path, s := range status {
		switch {
		case s.Worktree == git.Untracked:
			untracked = append(untracked, path)
		case s.Staging != git.Unmodified:
			staged[path] = getStatus(s.Staging)
		case s.Worktree != git.Unmodified:
			modified[path] = getStatus(s.Worktree)
		}
	}
	return
}
