package core

import (
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type KarnaLambdas struct {
	Client *lambda.Client
}

type KarnaAPIGateway struct {
	Client *apigateway.Client
}

type KarnaEC2 struct {
	Client *ec2.Client
}

type KarnaS3 struct {
	Client *s3.Client
}
type KarnaEC2Model struct {
	Instances      []ec2.Instance
	SecurityGroups []ec2.SecurityGroup
	Subnets        []ec2.Subnet
	VPCS           []string
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

type KarnaDeployment struct {
	Src          string            `json:"src"`
	Key          string            `json:"key"`
	File         string            `json:"file"`
	FunctionName string            `json:"functionName"`
	Aliases      map[string]string `json:"aliases"`
	Bucket       string            `json:"bucket"`
}
type KarnaConfigFile struct {
	Global      map[string]string `json:"global"`
	Deployments []KarnaDeployment `json:"deployments"`
	Path        string
}
