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
