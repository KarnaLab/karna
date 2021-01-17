package deploy

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/apigateway"
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

	if _, err = lambdaModel.getFunctionByFunctionName(targetDeployment.FunctionName); err != nil {
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
		if err = s3Model.upload(targetDeployment, output); err != nil {
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

	if err = lambdaModel.syncAlias(targetDeployment, *alias); err != nil {
		return timeElapsed, err
	}

	logger.Log("Done")

	if (targetDeployment.Prune.Alias) || (targetDeployment.Prune.Keep > 0) {
		logger.Log("Prune versions...")

		if err = lambdaModel.prune(targetDeployment); err != nil {
			return timeElapsed, err
		}
		logger.Log("Done")
	}

	if len(targetDeployment.API.ID) > 0 {
		logger.Log("Deploy to API Gateway...")

		_, err := agwModel.getRESTAPI(targetDeployment.API.ID)

		if err != nil {
			return timeElapsed, err
		}

		var currentResource string
		var parentResource apigateway.Resource

		resources, err := agwModel.getResources(targetDeployment.API.ID)

		if err != nil {
			return timeElapsed, err
		}

		for _, resource := range resources.Items {
			if *resource.Path == "/"+targetDeployment.API.Resource {
				currentResource = *resource.Id

				if _, ok := resource.ResourceMethods[targetDeployment.API.HTTPMethod]; !ok {
					logger.Log("Method do not exists for this resource, try to create it...")

					if _, err = agwModel.putMethod(targetDeployment.API.ID, currentResource, targetDeployment.API.HTTPMethod); err != nil {
						return timeElapsed, err
					}

					logger.Log("Method created!")
				}
			}
			// TODO: Make it able to transverse all Resources tree.
			if *resource.Path == "/" {
				parentResource = resource
			}
		}

		if currentResource == "" {
			logger.Log("Resource do not exists, try to create it...")
			resource, err := agwModel.createResource(targetDeployment.API.ID, *parentResource.Id, targetDeployment.API.Resource)

			currentResource = *resource.Id

			if err != nil {
				return timeElapsed, err
			}
			logger.Log("Resource created!")

			logger.Log("Try to create method into resource...")
			_, err = agwModel.putMethod(targetDeployment.API.ID, *resource.Id, targetDeployment.API.HTTPMethod)

			if err != nil {
				return timeElapsed, err
			}
			logger.Log("Method created!")
		}

		var integrationURI string

		integration, notFound, err := agwModel.getIntegration(targetDeployment.API.ID, currentResource, targetDeployment.API.HTTPMethod)

		if err != nil {
			if notFound {
				logger.Log("Try to create an integration for the method...")

				newIntegration, err := agwModel.putIntegration(targetDeployment.API.ID, currentResource, targetDeployment.API.HTTPMethod, targetDeployment.FunctionName)

				if err != nil {
					return timeElapsed, err
				}

				integrationURI = *newIntegration.Uri

				logger.Log("Integration created!")
			} else {
				return timeElapsed, err
			}
		} else {
			integrationURI = *integration.Uri
		}

		index := strings.Index(integrationURI, "${stageVariables.lambdaAlias}")

		if index == -1 {
			return timeElapsed, fmt.Errorf("Integration method is not valid. Must specify <function-name>:${stageVariables.lambdaAlias}")
		}

		stage, notFound, err := agwModel.getStage(targetDeployment.API.ID, *alias)

		if err != nil {
			if notFound {

				_, err := agwModel.createDeployment(targetDeployment.API.ID, *alias)

				if err != nil {
					return timeElapsed, err
				}

				if _, err = lambdaModel.addPermission(targetDeployment.FunctionName + ":" + *alias); err != nil {
					return timeElapsed, err
				}

				stage, _, err = agwModel.getStage(targetDeployment.API.ID, *alias)

				if err != nil {
					return timeElapsed, err
				}

			} else {
				return timeElapsed, err
			}
		}

		if stage.Variables["lambdaModelAlias"] == "" || stage.Variables["lamdaAlias"] != *alias {
			_, err := agwModel.updateStage(targetDeployment.API.ID, *alias)

			if err != nil {
				return timeElapsed, err
			}
		}

		logger.Log("API available at: https://" + targetDeployment.API.ID + ".execute-api." + agwModel.Client.Region + ".amazonaws.com/" + *alias + "/" + targetDeployment.API.Resource)

		logger.Log("Done")
	}

	timeElapsed = time.Since(startTime).String()
	return
}
