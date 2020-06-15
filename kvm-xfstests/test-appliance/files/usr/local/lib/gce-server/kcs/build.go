package main

import (
	"log"
	"os/exec"

	"example.com/gce-server/util"
)

const (
	buildScriptDir = "/usr/local/lib/gce-build-kernel"
)

func buildKernel(c LTMRequest) LTMRespond {
	repo := util.CloneCommit(c.Options.GitRepo, c.Options.CommitID)
	log.Printf("get build request: %+v\n", repo)

	cmd := exec.Command("bash", "-x", buildScriptDir, repo.GetDir(), util.GetConfig()["GS_BUCKET"])
	util.CheckRun(cmd, repo.GetDir())
	respond := LTMRespond{true}
	return respond
}
