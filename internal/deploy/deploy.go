package deploy

import (
	"time"

	"github.com/karnalab/karna/core"
)

func Run(target, alias *string) (timeElapsed string, err error) {
	var logger *core.KarnaLogger
	startTime := time.Now()

	configFile, err := getConfigFile()

	if err != nil {
		return timeElapsed, err
	}

	targetDeployment := getTargetDeployment(configFile, target)

	logger.Log("Checking requirements...")

	checkRequirements(targetDeployment, *alias)

	logger.Log("Done")

	var source string

	if targetDeployment.Executable == "" {
		source = configFile.Path + "/" + targetDeployment.Src
	} else {
		source = configFile.Path + "/" + targetDeployment.Src + "/" + targetDeployment.Executable
	}

	var outputPathWithoutArchive = configFile.Path + "/.karna/" + targetDeployment.FunctionName + "/" + *alias
	var output = configFile.Path + "/.karna/" + targetDeployment.FunctionName + "/" + *alias + "/" + targetDeployment.File

	logger.Log("Building archive...")

	if err = zipArchive(source, output, outputPathWithoutArchive, len(targetDeployment.Executable) > 0); err != nil {
		return timeElapsed, err
	}

	if targetDeployment.Bucket != "" {
		if err = core.S3.Upload(targetDeployment, output); err != nil {
			return timeElapsed, err
		}
	}

	logger.Log("Done")
	logger.Log("Updating function code...")

	err = core.Lambda.UpdateFunctionCode(targetDeployment, output)

	if err != nil {
		return timeElapsed, err
	}

	logger.Log("Done")
	logger.Log("Publishing function...")

	if err = core.Lambda.PublishFunction(targetDeployment); err != nil {
		return timeElapsed, err
	}

	logger.Log("Done")

	if err = core.Lambda.SyncAlias(targetDeployment, *alias); err != nil {
		return timeElapsed, err
	}

	logger.Log("Done")

	if (targetDeployment.Prune.Alias) || (targetDeployment.Prune.Keep > 0) {
		if err = core.Lambda.Prune(targetDeployment); err != nil {
			return timeElapsed, err
		}
	}

	timeElapsed = time.Since(startTime).String()
	return
}
