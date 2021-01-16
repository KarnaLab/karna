package deploy

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

var logger *KarnaLogger

type KarnaAPIGatewayModel struct {
	Client *apigateway.Client
}

func (karnaAGW *KarnaAPIGatewayModel) init() (err error) {
	cfg, err := external.LoadDefaultAWSConfig()

	if err != nil {
		return fmt.Errorf("unable to load SDK config, " + err.Error())
	}

	karnaAGW.Client = apigateway.New(cfg)
	return
}

func (karnaAGW *KarnaAPIGatewayModel) GetIntegration(APIID, resourceID, httpMethod string) (result *apigateway.GetIntegrationResponse, err error) {
	input := &apigateway.GetIntegrationInput{RestApiId: aws.String(APIID), ResourceId: aws.String(resourceID), HttpMethod: aws.String(httpMethod)}

	req := karnaAGW.Client.GetIntegrationRequest(input)

	result, err = req.Send(context.Background())

	return
}

func (karnaAGW *KarnaAPIGatewayModel) GetStage(APIID, stage string) (result *apigateway.GetStageResponse, notFound bool, err error) {
	input := &apigateway.GetStageInput{RestApiId: aws.String(APIID), StageName: aws.String(stage)}

	req := karnaAGW.Client.GetStageRequest(input)

	result, err = req.Send(context.Background())

	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case apigateway.ErrCodeNotFoundException:
			notFound = true
			break
		}
	}
	return
}

func (karnaAGW *KarnaAPIGatewayModel) UpdateStage(APIID, stage string) (result *apigateway.UpdateStageResponse, err error) {
	var updatedValues = apigateway.PatchOperation{
		Op:    apigateway.OpReplace,
		Path:  aws.String("/variables/lambdaAlias"),
		Value: aws.String(stage),
	}
	input := &apigateway.UpdateStageInput{RestApiId: aws.String(APIID), StageName: aws.String(stage), PatchOperations: []apigateway.PatchOperation{updatedValues}}

	req := karnaAGW.Client.UpdateStageRequest(input)

	result, err = req.Send(context.Background())

	return
}

func (karnaAGW *KarnaAPIGatewayModel) CreateStage(APIID, stageName, deploymentID string) (result *apigateway.CreateStageResponse, err error) {
	input := &apigateway.CreateStageInput{
		DeploymentId: aws.String(deploymentID),
		RestApiId:    aws.String(APIID),
		StageName:    aws.String(stageName),
		Description:  aws.String("Stage for " + stageName + " alias created by Karna"),
		Variables: map[string]string{
			"lambdaAlias": stageName,
		},
	}

	req := karnaAGW.Client.CreateStageRequest(input)

	result, err = req.Send(context.Background())

	return
}

func (karnaAGW *KarnaAPIGatewayModel) CreateDeployment(APIID, stageName string) (result *apigateway.CreateDeploymentResponse, err error) {
	input := &apigateway.CreateDeploymentInput{
		RestApiId:   aws.String(APIID),
		StageName:   aws.String(stageName),
		Description: aws.String("Deployment for stage " + stageName + " created by Karna"),
		Variables: map[string]string{
			"lambdaAlias": stageName,
		},
	}

	req := karnaAGW.Client.CreateDeploymentRequest(input)

	result, err = req.Send(context.Background())

	return
}

func (karnaAGW *KarnaAPIGatewayModel) GetRESTAPI(APIID string) (result *apigateway.GetRestApiResponse, err error) {
	input := &apigateway.GetRestApiInput{
		RestApiId: aws.String(APIID),
	}

	req := karnaAGW.Client.GetRestApiRequest(input)

	result, err = req.Send(context.Background())

	return
}

func (karnaAGW *KarnaAPIGatewayModel) GetResource(APIID, resourceID string) (result *apigateway.GetResourceResponse, err error) {
	input := &apigateway.GetResourceInput{
		RestApiId:  aws.String(APIID),
		ResourceId: aws.String(resourceID),
	}

	req := karnaAGW.Client.GetResourceRequest(input)

	result, err = req.Send(context.Background())

	return
}
