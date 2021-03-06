package deploy

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

var logger *KarnaLogger

type KarnaAPIGatewayModel struct {
	Client *apigateway.Client
	Sts    *sts.Client
}

func (karnaAGW *KarnaAPIGatewayModel) init() (err error) {
	cfg, err := external.LoadDefaultAWSConfig()

	if err != nil {
		return fmt.Errorf("unable to load SDK config, " + err.Error())
	}

	karnaAGW.Client = apigateway.New(cfg)

	karnaAGW.Sts = sts.New(cfg)

	return
}

func (karnaAGW *KarnaAPIGatewayModel) getIntegration(APIID, resourceID, httpMethod string) (result *apigateway.GetIntegrationResponse, notFound bool, err error) {
	input := &apigateway.GetIntegrationInput{
		RestApiId:  aws.String(APIID),
		ResourceId: aws.String(resourceID),
		HttpMethod: aws.String(httpMethod),
	}

	req := karnaAGW.Client.GetIntegrationRequest(input)

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

func (karnaAGW *KarnaAPIGatewayModel) getStage(APIID, stage string) (result *apigateway.GetStageResponse, notFound bool, err error) {
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

func (karnaAGW *KarnaAPIGatewayModel) updateStage(APIID, stage string) (result *apigateway.UpdateStageResponse, err error) {
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

func (karnaAGW *KarnaAPIGatewayModel) createStage(APIID, stageName, deploymentID string) (result *apigateway.CreateStageResponse, err error) {
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

func (karnaAGW *KarnaAPIGatewayModel) createDeployment(APIID, stageName string) (result *apigateway.CreateDeploymentResponse, err error) {
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

func (karnaAGW *KarnaAPIGatewayModel) getRESTAPI(APIID string) (result *apigateway.GetRestApiResponse, err error) {
	input := &apigateway.GetRestApiInput{
		RestApiId: aws.String(APIID),
	}

	req := karnaAGW.Client.GetRestApiRequest(input)

	result, err = req.Send(context.Background())

	return
}

func (karnaAGW *KarnaAPIGatewayModel) getResource(APIID, resourceID string) (result *apigateway.GetResourceResponse, err error) {
	input := &apigateway.GetResourceInput{
		RestApiId:  aws.String(APIID),
		ResourceId: aws.String(resourceID),
	}

	req := karnaAGW.Client.GetResourceRequest(input)

	result, err = req.Send(context.Background())

	return
}

func (karnaAGW *KarnaAPIGatewayModel) getResources(APIID string) (result *apigateway.GetResourcesResponse, err error) {
	input := &apigateway.GetResourcesInput{
		RestApiId: aws.String(APIID),
	}

	req := karnaAGW.Client.GetResourcesRequest(input)

	result, err = req.Send(context.Background())

	return
}

func (karnaAGW *KarnaAPIGatewayModel) createResource(APIID, parentID, pathPart string) (result *apigateway.CreateResourceResponse, err error) {
	input := &apigateway.CreateResourceInput{
		RestApiId: aws.String(APIID),
		ParentId:  aws.String(parentID),
		PathPart:  aws.String(pathPart),
	}

	req := karnaAGW.Client.CreateResourceRequest(input)

	result, err = req.Send(context.Background())

	return
}

func (karnaAGW *KarnaAPIGatewayModel) putMethod(APIID, resourceID, httpMethod string) (result *apigateway.PutMethodResponse, err error) {
	input := &apigateway.PutMethodInput{
		RestApiId:         aws.String(APIID),
		ResourceId:        aws.String(resourceID),
		HttpMethod:        aws.String(httpMethod),
		AuthorizationType: aws.String("NONE"),
	}

	req := karnaAGW.Client.PutMethodRequest(input)

	result, err = req.Send(context.Background())

	return
}

func (karnaAGW *KarnaAPIGatewayModel) getAccountID() (result *sts.GetAccessKeyInfoResponse, err error) {
	credentials, err := karnaAGW.Client.Credentials.Retrieve()
	stsInput := &sts.GetAccessKeyInfoInput{
		AccessKeyId: aws.String(credentials.AccessKeyID),
	}

	req := karnaAGW.Sts.GetAccessKeyInfoRequest(stsInput)

	result, err = req.Send(context.Background())

	return
}

func (karnaAGW *KarnaAPIGatewayModel) putIntegration(APIID, resourceID, httpMethod, functionName string) (result *apigateway.PutIntegrationResponse, err error) {
	region := karnaAGW.Client.Region

	account, err := karnaAGW.getAccountID()

	if err != nil {
		return
	}

	input := &apigateway.PutIntegrationInput{
		RestApiId:             aws.String(APIID),
		ResourceId:            aws.String(resourceID),
		HttpMethod:            aws.String(httpMethod),
		IntegrationHttpMethod: aws.String("POST"),
		Type:                  apigateway.IntegrationTypeAwsProxy,
		Uri:                   aws.String("arn:aws:apigateway:" + region + ":lambda:path/2015-03-31/functions/arn:aws:lambda:" + region + ":" + *account.Account + ":function:" + functionName + ":${stageVariables.lambdaAlias}/invocations"),
	}

	req := karnaAGW.Client.PutIntegrationRequest(input)

	result, err = req.Send(context.Background())

	if err != nil {
		return
	}

	inputResponse := &apigateway.PutIntegrationResponseInput{
		RestApiId:         aws.String(APIID),
		ResourceId:        aws.String(resourceID),
		HttpMethod:        aws.String(httpMethod),
		StatusCode:        aws.String("200"),
		ResponseTemplates: map[string]string{"application/json": ""},
	}

	integrationReq := karnaAGW.Client.PutIntegrationResponseRequest(inputResponse)

	_, err = integrationReq.Send(context.Background())

	return
}

func (karnaAGW *KarnaAPIGatewayModel) deleteStage(APIID, stageName string) (result *apigateway.DeleteStageResponse, err error) {
	input := &apigateway.DeleteStageInput{
		RestApiId: aws.String(APIID),
		StageName: aws.String(stageName),
	}

	req := karnaAGW.Client.DeleteStageRequest(input)

	result, err = req.Send(context.Background())

	return
}
