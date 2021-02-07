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

func RemoveAlias(functionName, alias *string) (timeElapsed string, err error) {
	var logger KarnaLogger

	startTime := time.Now()

	configFile, err := getConfigFile()

	if err != nil {
		return timeElapsed, err
	}

	targetDeployment, err := getTargetDeployment(configFile, *functionName)

	if err != nil {
		return timeElapsed, err
	}

	logger.Log("Checking if alias exists...")

	aliases, _ := lambdaModel.getAliasesByFunctionName(*functionName)

	if a := findAlias(aliases, *alias); a == nil {
		return timeElapsed, fmt.Errorf("Alias do not exists, operation aborted")
	}

	logger.Log("Done")

	logger.Log("Removing in progress...")

	if err := lambdaModel.deleteAlias(*functionName, *alias, targetDeployment); err != nil {
		return timeElapsed, err
	}

	logger.Log("Alias removed")

	if len(targetDeployment.API.ID) > 0 {
		if _, err := agwModel.deleteStage(targetDeployment.API.ID, *alias); err != nil {
			return timeElapsed, err
		}

		logger.Log("Stage removed")
	}

	logger.Log("Done")

	timeElapsed = time.Since(startTime).String()
	return
}

func Deploy(functionName, alias *string) (timeElapsed string, err error) {
	var logger KarnaLogger

	startTime := time.Now()

	configFile, err := getConfigFile()

	if err != nil {
		return timeElapsed, err
	}

	targetDeployment, err := getTargetDeployment(configFile, *functionName)

	if err != nil {
		return timeElapsed, err
	}

	logger.Log("Checking requirements...")

	if err = checkRequirements(targetDeployment, *alias); err != nil {
		return timeElapsed, err
	}

	logger.Log("Done")

	if _, err = lambdaModel.getFunctionByFunctionName(*functionName); err != nil {
		return timeElapsed, err
	}

	var source string

	if targetDeployment.Executable == "" {
		source = configFile.Path + "/" + targetDeployment.Src
	} else {
		source = configFile.Path + "/" + targetDeployment.Src + "/" + targetDeployment.Executable
	}

	var outputPathWithoutArchive = configFile.Path + "/.karna/" + *functionName + "/" + *alias
	var output = configFile.Path + "/.karna/" + *functionName + "/" + *alias + "/" + targetDeployment.File

	logger.Log("Building archive in progress...")

	if err = zipArchive(source, output, outputPathWithoutArchive, len(targetDeployment.Executable) > 0); err != nil {
		return timeElapsed, err
	}

	if targetDeployment.Bucket != "" {
		if err = s3Model.upload(targetDeployment, output); err != nil {
			return timeElapsed, err
		}
	}

	logger.Log("Done")
	logger.Log("Updating function code in progress...")

	err = lambdaModel.UpdateFunctionCode(targetDeployment, output, *functionName)

	if err != nil {
		return timeElapsed, err
	}

	logger.Log("Done")
	logger.Log("Publishing function in progress...")

	if err = lambdaModel.PublishFunction(*functionName, targetDeployment); err != nil {
		return timeElapsed, err
	}

	logger.Log("Done")

	logger.Log("Synchronizing alias in progress...")

	if err = lambdaModel.syncAlias(targetDeployment, *alias, *functionName); err != nil {
		return timeElapsed, err
	}

	logger.Log("Done")

	if ok := targetDeployment.Versions.Keep > 0; ok {
		logger.Log("Versions removing in progress...")

		if err = lambdaModel.removeVersions(*functionName, targetDeployment); err != nil {
			return timeElapsed, err
		}
		logger.Log("Done")
	}

	if len(targetDeployment.API.ID) > 0 {
		logger.Log("Deployment to API Gateway in progress...")

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
					logger.Log("Method do not exists for this resource, creation in progress...")

					if _, err = agwModel.putMethod(targetDeployment.API.ID, currentResource, targetDeployment.API.HTTPMethod); err != nil {
						return timeElapsed, err
					}

					logger.Log("Done")
				}
			}
			// TODO: Make it able to transverse all Resources tree.
			if *resource.Path == "/" {
				parentResource = resource
			}
		}

		if currentResource == "" {
			logger.Log("Resource do not exists, creation in progress...")
			resource, err := agwModel.createResource(targetDeployment.API.ID, *parentResource.Id, targetDeployment.API.Resource)

			currentResource = *resource.Id

			if err != nil {
				return timeElapsed, err
			}
			logger.Log("Done")

			logger.Log("Method creation in progress...")
			_, err = agwModel.putMethod(targetDeployment.API.ID, *resource.Id, targetDeployment.API.HTTPMethod)

			if err != nil {
				return timeElapsed, err
			}
			logger.Log("Done")
		}

		var integrationURI string

		integration, notFound, err := agwModel.getIntegration(targetDeployment.API.ID, currentResource, targetDeployment.API.HTTPMethod)

		if err != nil {
			if notFound {
				logger.Log("Try to create an integration for the method...")

				newIntegration, err := agwModel.putIntegration(targetDeployment.API.ID, currentResource, targetDeployment.API.HTTPMethod, *functionName)

				if err != nil {
					return timeElapsed, err
				}

				integrationURI = *newIntegration.Uri

				logger.Log("Done")
			} else {
				return timeElapsed, err
			}
		} else {
			integrationURI = *integration.Uri
		}

		index := strings.Index(integrationURI, "${stageVariables.lambdaAlias}")

		if index == -1 {
			return timeElapsed, fmt.Errorf("Integration method is not valid. You must specify <function-name>:${stageVariables.lambdaAlias}")
		}

		stage, notFound, err := agwModel.getStage(targetDeployment.API.ID, *alias)

		if err != nil {
			if notFound {

				_, err := agwModel.createDeployment(targetDeployment.API.ID, *alias)

				if err != nil {
					return timeElapsed, err
				}

				if _, err = lambdaModel.addPermission(*functionName + ":" + *alias); err != nil {
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

		logger.Log("Done")
		logger.Log("API available at: https://" + targetDeployment.API.ID + ".execute-api." + agwModel.Client.Region + ".amazonaws.com/" + *alias + "/" + targetDeployment.API.Resource)
	}

	timeElapsed = time.Since(startTime).String()
	return
}
