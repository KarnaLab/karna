package deploy

import (
	"fmt"
	"strings"
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

	targetDeployment, err := getTargetDeployment(configFile, target)

	if err != nil {
		return timeElapsed, err
	}

	logger.Log("Checking requirements...")

	if err = checkRequirements(targetDeployment, *alias); err != nil {
		return timeElapsed, err
	}

	logger.Log("Done")

	if _, err = core.Lambda.GetFunctionByFunctionName(targetDeployment.FunctionName); err != nil {
		return timeElapsed, err
	}

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

	logger.Log("Synchronize alias...")

	if err = core.Lambda.SyncAlias(targetDeployment, *alias); err != nil {
		return timeElapsed, err
	}

	logger.Log("Done")

	if (targetDeployment.Prune.Alias) || (targetDeployment.Prune.Keep > 0) {
		logger.Log("Prune versions...")

		if err = core.Lambda.Prune(targetDeployment); err != nil {
			return timeElapsed, err
		}
		logger.Log("Done")
	}

	if len(targetDeployment.API.ID) > 0 {
		logger.Log("Deploy to API Gateway...")
		var shouldReDeploy bool

		apisTree := core.AGW.BuildAGWTree()

		var currentAPI core.KarnaAGWAPI
		var currentResource map[string]interface{}

		for _, api := range apisTree {
			if *api.API.Id == targetDeployment.API.ID {
				currentAPI = api
			}
		}

		if currentAPI.API.Name == nil {
			return timeElapsed, fmt.Errorf("API not found")
		}

		for _, resource := range currentAPI.Resources {
			if resource["Id"] == targetDeployment.API.Resource {
				currentResource = resource
			}
		}

		if currentResource["Id"] == nil {
			return timeElapsed, fmt.Errorf("Resource not found")
		}

		integration, err := core.AGW.GetIntegration(targetDeployment.API.ID, targetDeployment.API.Resource, targetDeployment.API.HTTPMethod)

		if err != nil {
			return timeElapsed, err
		}

		index := strings.Index(*integration.Uri, "${stageVariables.lambdaAlias}")

		if index == -1 {
			return timeElapsed, fmt.Errorf("Integration method is not valid. Must specify ${stageVariable.lambdaAlias}")
		}

		stage, notFound, err := core.AGW.GetStage(targetDeployment.API.ID, *alias)

		if err != nil {
			if notFound {
				shouldReDeploy = true
				if _, err = core.AGW.CreateStage(targetDeployment.API.ID, *alias, "1tbqsq"); err != nil {
					return timeElapsed, err
				}
			}

			return timeElapsed, err
		}

		if stage.Variables["lambdaAlias"] == "" || stage.Variables["lamdaAlias"] != *alias {
			updated, err := core.AGW.UpdateStage(targetDeployment.API.ID, *alias)

			if err != nil {
				return timeElapsed, err
			}
			fmt.Println(updated)
		}

		if shouldReDeploy {
			// Redeploy API
		}

		logger.Log("Done")
	}

	timeElapsed = time.Since(startTime).String()
	return
}
