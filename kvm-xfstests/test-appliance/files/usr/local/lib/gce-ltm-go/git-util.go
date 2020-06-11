package main

import (
	"log"
	"os"
	"os/exec"
	"time"
)

const (
	rootdir       = "/root/repositories"
	checkInterval = 5
)

// clone a repository into specified directory
// remove directory if already exists
// only public repos are supported
func gitClone(url string, dirName string) {

	errDir := os.MkdirAll(rootdir, 755)
	if errDir != nil {
		log.Fatal(errDir)
	}

	dir := rootdir + "/" + dirName
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

}

// get the newest commit id on a local branch
// the specified branch must exist locally
func gitCommit(dirName string, branch string) string {
	dir := rootdir + "/" + dirName
	stat, err := os.Stat(dir)
	if err != nil || !stat.IsDir() {
		log.Fatalf("directory %s does not exist!", dir)
	}
	cmd := exec.Command("git", "checkout", branch)
	checkRun(cmd, dir)
	cmd = exec.Command("git", "rev-parse", "@")
	commit := checkOutput(cmd, dir)
	return commit[:len(commit)-1]
}

// pull the newest code from upstream
// return the newest commit id
func gitPull(dirName string, branch string) string {
	dir := rootdir + "/" + dirName
	stat, err := os.Stat(dir)
	if err != nil || !stat.IsDir() {
		log.Fatalf("directory %s does not exist!", dir)
	}
	cmd := exec.Command("git", "checkout", branch)
	checkRun(cmd, dir)
	cmd = exec.Command("git", "pull")
	checkRun(cmd, dir)
	cmd = exec.Command("git", "rev-parse", "@")
	commit := checkOutput(cmd, dir)
	return commit[:len(commit)-1]
}

// watch a specified branch in a repo
// print newest commit id when upstream changes
func gitWatch(url string, dirName string, branch string) {
	gitClone(url, dirName)
	currCommit := gitCommit(dirName, branch)
	newCommit := currCommit
	for true {
		time.Sleep(checkInterval * time.Second)
		newCommit = gitPull(dirName, branch)
		if newCommit != currCommit {
			log.Println("new commit detected")
			currCommit = newCommit
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
