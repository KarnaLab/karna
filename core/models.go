package core

import (
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

type KarnaLambdas struct {
	Client *lambda.Client
}

type KarnaAPIGateway struct {
	Client *apigateway.Client
}

type KarnaLambda struct {
	FunctionConfiguration lambda.FunctionConfiguration
	Layers                []lambda.Layer
	VPC                   string
	Versions              []lambda.FunctionConfiguration
	Policy                map[string][]string
}

type awsPolicyStatementCondition struct {
	ArnLike map[string]string
}
type awsPolicyStatementPrincipal struct {
	Service string
}
type awsPolicyStatement struct {
	Action    string
	Effect    string
	Resource  string
	ID        string `json:"$id"`
	Condition awsPolicyStatementCondition
	Principal awsPolicyStatementPrincipal
}

type awsPolicy struct {
	Version   string
	ID        string
	Statement []awsPolicyStatement
}

type KarnaAGWStage struct {
	Name         string
	Stage        string
	Uuid         string
	Distribution string
}

type KarnaAGWAPI struct {
	API       apigateway.RestApi
	Resources []map[string]interface{}
	Stages    []KarnaAGWStage
}
