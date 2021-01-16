package deploy

import (
	"fmt"
	"strings"
	"time"
)

var lambdaModel KarnaLambdaModel
var agwModel KarnaAPIGatewayModel
var s3Model KarnaS3Model

func init() {
	lambdaModel.init()
	agwModel.init()
	s3Model.init()
}

func Run(target, alias *string) (timeElapsed string, err error) {
	var logger KarnaLogger

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

	if _, err = lambdaModel.GetFunctionByFunctionName(targetDeployment.FunctionName); err != nil {
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
		if err = s3Model.Upload(targetDeployment, output); err != nil {
			return timeElapsed, err
		}
	}

	logger.Log("Done")
	logger.Log("Updating function code...")

	err = lambdaModel.UpdateFunctionCode(targetDeployment, output)

	if err != nil {
		return timeElapsed, err
	}

	logger.Log("Done")
	logger.Log("Publishing function...")

	if err = lambdaModel.PublishFunction(targetDeployment); err != nil {
		return timeElapsed, err
	}

	logger.Log("Done")

	logger.Log("Synchronize alias...")

	if err = lambdaModel.SyncAlias(targetDeployment, *alias); err != nil {
		return timeElapsed, err
	}

	logger.Log("Done")

	if (targetDeployment.Prune.Alias) || (targetDeployment.Prune.Keep > 0) {
		logger.Log("Prune versions...")

		if err = lambdaModel.Prune(targetDeployment); err != nil {
			return timeElapsed, err
		}
		logger.Log("Done")
	}

	if len(targetDeployment.API.ID) > 0 {
		logger.Log("Deploy to API Gateway...")

		var currentResource map[string]interface{}

		_, err := agwModel.GetRESTAPI(targetDeployment.API.ID)

		if err != nil {
			return timeElapsed, err
		}

		_, err = agwModel.GetResource(targetDeployment.API.ID, targetDeployment.API.Resource)

		if err != nil {
			return timeElapsed, err
		}

		integration, err := agwModel.GetIntegration(targetDeployment.API.ID, targetDeployment.API.Resource, targetDeployment.API.HTTPMethod)

		if err != nil {
			return timeElapsed, err
		}

		index := strings.Index(*integration.Uri, "${stageVariables.lambdaModelAlias}")

		if index == -1 {
			return timeElapsed, fmt.Errorf("Integration method is not valid. Must specify ${stageVariable.lambdaModelAlias}")
		}

		stage, notFound, err := agwModel.GetStage(targetDeployment.API.ID, *alias)

		if err != nil {
			if notFound {

				_, err := agwModel.CreateDeployment(targetDeployment.API.ID, *alias)

				if err != nil {
					return timeElapsed, err
				}

				if _, err = lambdaModel.AddPermission(targetDeployment.FunctionName, *alias); err != nil {
					return timeElapsed, err
				}

				stage, _, err = agwModel.GetStage(targetDeployment.API.ID, *alias)

				if err != nil {
					return timeElapsed, err
				}

			} else {
				return timeElapsed, err
			}
		}

		if stage.Variables["lambdaModelAlias"] == "" || stage.Variables["lamdaAlias"] != *alias {
			_, err := agwModel.UpdateStage(targetDeployment.API.ID, *alias)

			if err != nil {
				return timeElapsed, err
			}
		}

		logger.Log("API available at: https://" + targetDeployment.API.ID + ".execute-api." + agwModel.Client.Region + ".amazonaws.com/" + *alias + currentResource["Path"].(string))

		logger.Log("Done")
	}

	timeElapsed = time.Since(startTime).String()
	return
}
