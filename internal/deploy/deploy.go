package deploy

import (
	"fmt"
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

	targetFunction, err := getTargetFunction(config, target, alias)
	fmt.Println(targetFunction)
	if err != nil {
		return timeElapsed, err
	}

	var input = config.Path + "/" + targetFunction.Input

	if _, err := os.Stat(input); os.IsNotExist(err) {
		return timeElapsed, err
	}

	var output = config.Path + "/.karna/" + targetFunction.Name + "/" + targetFunction.Output

	logger.Log("Building archive...")

	if err = zipArchive(input, output); err != nil {
		return timeElapsed, err
	}

	if targetFunction.S3.Bucket != "" {
		if err = core.S3.Upload(targetFunction, output); err != nil {
			return timeElapsed, err
		}
	}

	logger.Log("Done")
	logger.Log("Updating function code...")

	err = core.Lambda.UpdateFunctionCode(targetFunction, output)

	if err != nil {
		return timeElapsed, err
	}

	logger.Log("Done")
	logger.Log("Publishing function...")

	if err = core.Lambda.PublishFunction(targetFunction); err != nil {
		return timeElapsed, err
	}

	logger.Log("Done")

	if err = core.Lambda.SyncAlias(targetFunction, *alias); err != nil {
		return timeElapsed, err
	}

	logger.Log("Done")

	if (targetFunction.Prune.Alias) || (targetFunction.Prune.Keep > 0) {
		if err = core.Lambda.Prune(targetFunction); err != nil {
			return timeElapsed, err
		}
	}

	timeElapsed = time.Since(startTime).String()

	return
}
