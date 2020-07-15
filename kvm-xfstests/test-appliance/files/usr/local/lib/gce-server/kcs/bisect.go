package main

import (
	"fmt"

	"gce-server/logging"
	"gce-server/server"
	"gce-server/util"
)

// StartBisect starts a git bisect task.
func StartBisect(c server.TaskRequest, testID string) {
	log := server.Log.WithField("testID", testID)
	log.Info("Start git bisect")

	bisectLog := logging.KCSLogDir + testID + ".log"
	subject := "xfstests KCS bisect failure " + testID
	defer util.ReportFailure(log, bisectLog, c.Options.ReportEmail, subject)



	config, err := util.GetConfig(util.GceConfigFile)
	logging.CheckPanic(err, log, "Failed to get config")

	gsBucket := config.Get("GS_BUCKET")
	gsPath := fmt.Sprintf("gs://%s/kernels/bzImage-%s-onerun", gsBucket, testID)

	runBuild(c.Options.GitRepo, c.Options.CommitID, gsBucket, gsPath, testID, buildLog, log)
}
