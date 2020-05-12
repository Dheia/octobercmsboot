package octobercmsboot

import (
	"fmt"
	gitpkg "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"os"
)

type git struct {
	Repository
	ReferenceName plumbing.ReferenceName
}

type Repository struct {
	URL, Path, Branch string
}

func newGit(repo Repository) *git {
	if len(repo.Branch) == 0 {
		repo.Branch = "master"
	}
	return &git{
		Repository:    repo,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", repo.Branch)),
	}
}

func (g git) install() {
	if _, err := os.Stat(g.Repository.Path); os.IsNotExist(err) {
		g.clone()
		return
	}
	if _, err := os.Stat(g.Repository.Path + "/.git"); os.IsNotExist(err) {
		g.init()
	}
	g.pull()
}

func (g git) clone() {
	_, err := gitpkg.PlainClone(g.Repository.Path, false, &gitpkg.CloneOptions{
		URL:           g.Repository.URL,
		ReferenceName: g.ReferenceName,
		Progress:      os.Stdout,
	})
	g.checkIfError(err)
}

func (g git) init() {
	currentRepo, err := gitpkg.PlainInit(g.Repository.Path, false)
	g.checkIfError(err)
	_, err = currentRepo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{g.Repository.URL},
	})
	g.checkIfError(err)
	err = currentRepo.Fetch(&gitpkg.FetchOptions{
		RemoteName: "origin",
	})
	g.checkIfError(err)
	err = currentRepo.CreateBranch(&config.Branch{
		Name:   g.Repository.Branch,
		Remote: "origin",
		Merge:  g.ReferenceName,
	})
	g.checkIfError(err)
}

func (g git) pull() {
	pullOptions := &gitpkg.PullOptions{
		RemoteName:    "origin",
		ReferenceName: g.ReferenceName,
	}
	r, err := gitpkg.PlainOpen(g.Repository.Path)
	g.checkIfError(err)
	worktree, err := r.Worktree()
	g.checkIfError(err)
	err = worktree.Pull(pullOptions)
	g.checkIfError(err)
}

func (g git) checkIfError(err error) {
	if err == nil {
		return
	}
	if err == gitpkg.NoErrAlreadyUpToDate {
		Info("%s => %s", g.Repository.URL, gitpkg.NoErrAlreadyUpToDate.Error())
		return
	}
	Error("url %s error: %s", g.Repository.URL, err)
	os.Exit(1)
}
