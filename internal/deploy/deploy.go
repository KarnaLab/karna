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

	core.LogSuccessMessage("Checking requirements...")

	checkRequirements(targetDeployment, *alias)

	core.LogSuccessMessage("Done")

	var source = configFile.Path + "/" + targetDeployment.Src
	var output = configFile.Path + "/.karna/" + targetDeployment.FunctionName + "/" + *alias + "/" + targetDeployment.File

	core.LogSuccessMessage("Building archive...")

	zipArchive(source, output)

	if targetDeployment.Bucket != "" {
		err := core.S3.Upload(targetDeployment, output)

		if err != nil {
			core.LogErrorMessage(err.Error())
			os.Exit(2)
		}
	}
	core.LogSuccessMessage("Done")

	core.LogSuccessMessage("Updating function code...")
	err := core.Lambda.UpdateFunctionCode(targetDeployment, output)

	if err != nil {
		core.LogErrorMessage(err.Error())
		os.Exit(2)
	}
	core.LogSuccessMessage("Done")

	core.LogSuccessMessage("Publishing function...")
	err = core.Lambda.PublishFunction(targetDeployment)

	if err != nil {
		core.LogErrorMessage(err.Error())
		os.Exit(2)
	}

	core.LogSuccessMessage("Done")

	err = core.Lambda.SyncAlias(targetDeployment, *alias)

	if err != nil {
		core.LogErrorMessage(err.Error())
		os.Exit(2)
	}

	core.LogSuccessMessage("Done")

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
