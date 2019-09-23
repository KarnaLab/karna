package deploy

import (
	"karna/core"
	"os"
	"time"
)

func Run(target *string, alias *string) (timeElapsed string) {
	startTime := time.Now()

	configFile := getConfigFile()
	targetDeployment := getTargetDeployment(configFile, target)

	checkRequirements(targetDeployment, *alias)

	var source = configFile.Path + "/" + targetDeployment.Src
	var output = configFile.Path + "/.karna/" + targetDeployment.File

	zipArchive(source, output)

	if targetDeployment.Bucket != "" {
		err := core.S3.Upload(targetDeployment, output)

		if err != nil {
			core.LogErrorMessage(err.Error())
			os.Exit(2)
		}
	}

	err := core.Lambda.UpdateFunctionCode(targetDeployment, output)

	if err != nil {
		core.LogErrorMessage(err.Error())
		os.Exit(2)
	}

	err = core.Lambda.PublishFunction(targetDeployment)

	if err != nil {
		core.LogErrorMessage(err.Error())
		os.Exit(2)
	}

	err = core.Lambda.SyncAlias(targetDeployment, *alias)

	if err != nil {
		core.LogErrorMessage(err.Error())
		os.Exit(2)
	}

	if (targetDeployment.Prune.Alias) || (targetDeployment.Prune.Keep > 0) {
		err := core.Lambda.Prune(targetDeployment)

		if err != nil {
			core.LogErrorMessage(err.Error())
			os.Exit(2)
		}
	}

	timeElapsed = time.Since(startTime).String()
	return
}
