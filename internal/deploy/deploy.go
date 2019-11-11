package deploy

import (
	"os"
	"time"

	"github.com/karnalab/karna/core"
)

func Run(target, alias *string) (timeElapsed string, err error) {
	var logger *core.KarnaLogger
	startTime := time.Now()

	config, err := core.GetConfig()

	if err != nil {
		return timeElapsed, err
	}

	targetDeployment, err := getTargetDeployment(config, target, alias)

	if err != nil {
		return timeElapsed, err
	}

	var input = config.Path + "/" + targetDeployment.Input

	if _, err := os.Stat(input); os.IsNotExist(err) {
		return timeElapsed, err
	}

	var output = config.Path + "/.karna/" + targetDeployment.Name + "/" + targetDeployment.Output

	logger.Log("Building archive...")

	if err = zipArchive(input, output); err != nil {
		return timeElapsed, err
	}

	if targetDeployment.S3.Bucket != "" {
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
