package util

import (
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/google/uuid"
)

const (
	rootdir       = "/root/repositories"
	checkInterval = 10
)

type Repository struct {
	url        string
	id         string
	branch     string
	currCommit string
	watching   bool
}

// clone a repository into a unique directory and checkout to branch
// return a Repository
// only public repos are supported
func CloneBranch(url string, branch string) Repository {
	r := Repository{"", "", "", "", false}
	r.url = url
	r.branch = branch
	id, _ := uuid.NewRandom()
	r.id = id.String()

	errDir := os.MkdirAll(rootdir, 755)
	if errDir != nil {
		log.Fatal(errDir)
	}

	dir := rootdir + "/" + r.id
	stat, err := os.Stat(dir)
	if err == nil {
		if stat.IsDir() {
			os.RemoveAll(dir)
		} else {
			os.Remove(dir)
		}
	}
	cmd := exec.Command("git", "clone", url, dir)
	CheckRun(cmd, rootdir)

	// checkout to branch and get HEAD commit id
	r.currCommit = r.GetCommit()
	return r
}

// clone a repository into a unique directory and checkout to commit
// return a Repository
// only public repos are supported
func CloneCommit(url string, commit string) Repository {
	r := Repository{"", "", "", "", false}
	r.url = url
	r.branch = ""
	r.currCommit = commit
	id, _ := uuid.NewRandom()
	r.id = id.String()

	errDir := os.MkdirAll(rootdir, 755)
	if errDir != nil {
		log.Fatal(errDir)
	}

	dir := rootdir + "/" + r.id
	stat, err := os.Stat(dir)
	if err == nil {
		if stat.IsDir() {
			os.RemoveAll(dir)
		} else {
			os.Remove(dir)
		}
	}
	cmd := exec.Command("git", "clone", url, dir)
	CheckRun(cmd, rootdir)

	// checkout to commit
	cmd = exec.Command("git", "checkout", r.currCommit)
	CheckRun(cmd, dir)
	return r
}

// get the newest commit id on a local branch
// the specified branch must exist locally
func (r *Repository) GetCommit() string {
	dir := rootdir + "/" + r.id
	stat, err := os.Stat(dir)
	if err != nil || !stat.IsDir() {
		log.Fatalf("directory %s does not exist!", dir)
	}
	cmd := exec.Command("git", "checkout", r.branch)
	CheckRun(cmd, dir)

	cmd = exec.Command("git", "rev-parse", "@")
	commit := CheckOutput(cmd, dir)
	return commit[:len(commit)-1]
}

// pull the newest code from upstream
func (r *Repository) Pull() {
	dir := rootdir + "/" + r.id
	stat, err := os.Stat(dir)
	if err != nil || !stat.IsDir() {
		log.Fatalf("directory %s does not exist!", dir)
	}
	cmd := exec.Command("git", "pull")
	CheckRun(cmd, dir)
}

// watch a specified branch in a repo
// print newest commit id when upstream changes
func (r *Repository) Watch() {
	if r.watching {
		return
	}
	r.watching = true
	for {
		time.Sleep(checkInterval * time.Second)
		r.Pull()
		newCommit := r.GetCommit()
		if newCommit != r.currCommit {
			r.currCommit = newCommit
			log.Println("new commit detected")
		}
	}
}

func (r *Repository) GetDir() string {
	return rootdir + "/" + r.id
}

func GetRepo(url string, id string, branch string, currCommit string, watching bool) Repository {
	return Repository{url, id, branch, currCommit, watching}
}
