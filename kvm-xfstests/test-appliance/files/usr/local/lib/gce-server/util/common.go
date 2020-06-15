package util

import (
	"log"
	"os"
	"os/exec"
)

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CheckRun(cmd *exec.Cmd, workDir string) {
	cmd.Dir = workDir
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("%s failed with error: %s\n", cmd.String(), err)
	}
}

func CheckOutput(cmd *exec.Cmd, workDir string) string {
	cmd.Dir = workDir
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("%s failed with error: %s\n", cmd.String(), err)
	}
	return string(out)
}
