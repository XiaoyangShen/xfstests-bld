package main

import (
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/google/uuid"
)

const (
	rootdir       = "/root/repositories"
	checkInterval = 5
)

type Repository struct {
	url        string
	id         string
	branch     string
	currCommit string
	watching   bool
}

// clone a repository into a unique directory
// remove directory if already exists
// should be called only once
// only public repos are supported
func (r *Repository) clone(url string, branch string) {
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
	checkRun(cmd, rootdir)

	r.currCommit = r.getCommit()
	r.watching = false
}

// get the newest commit id on a local branch
// the specified branch must exist locally
func (r *Repository) getCommit() string {
	dir := rootdir + "/" + r.id
	stat, err := os.Stat(dir)
	if err != nil || !stat.IsDir() {
		log.Fatalf("directory %s does not exist!", dir)
	}
	cmd := exec.Command("git", "checkout", r.branch)
	checkRun(cmd, dir)
	cmd = exec.Command("git", "rev-parse", "@")
	commit := checkOutput(cmd, dir)
	return commit[:len(commit)-1]
}

// pull the newest code from upstream
func (r *Repository) pull() {
	dir := rootdir + "/" + r.id
	stat, err := os.Stat(dir)
	if err != nil || !stat.IsDir() {
		log.Fatalf("directory %s does not exist!", dir)
	}
	cmd := exec.Command("git", "pull")
	checkRun(cmd, dir)
}

// watch a specified branch in a repo
// print newest commit id when upstream changes
func (r *Repository) watch() {
	if r.watching {
		return
	}
	r.watching = true
	for true {
		time.Sleep(checkInterval * time.Second)
		newCommit := r.getCommit()
		if newCommit != r.currCommit {
			log.Println("new commit detected")
			r.currCommit = newCommit
		}
	}
}

func checkRun(cmd *exec.Cmd, workDir string) {
	cmd.Dir = workDir
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("%s failed with error: %s\n", cmd.String(), err)
	}
}

func checkOutput(cmd *exec.Cmd, workDir string) string {
	cmd.Dir = workDir
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("%s failed with error: %s\n", cmd.String(), err)
	}
	return string(out)
}
