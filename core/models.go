package core

import (
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

//KarnaLambdaModel => Karna API for Lambda.
type KarnaLambdaModel struct {
	Client *lambda.Client
}

//KarnaAPIGatewayModel => Karna API for APIGateway.
type KarnaAPIGatewayModel struct {
	Client *apigateway.Client
}

//KarnaEC2Model => Karna API for EC2.
type KarnaEC2Model struct {
	Client *ec2.Client
}

//KarnaS3Model => Karna API for S3.
type KarnaS3Model struct {
	Client *s3.Client
}

//KarnaEC2 => Karna model for EC2.
type KarnaEC2 struct {
	Instances      []ec2.Instance
	SecurityGroups []ec2.SecurityGroup
	Subnets        []ec2.Subnet
	VPCS           []string
}

//KarnaLambda => Karna model for Lambda.
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

//KarnaAGWStage => Karna model for APIGateway Stage.
type KarnaAGWStage struct {
	Name         string
	Stage        string
	Uuid         string
	Distribution string
}

//KarnaAGWAPI => Karna model for APIGateway.
type KarnaAGWAPI struct {
	API       apigateway.RestApi
	Resources []map[string]interface{}
	Stages    []KarnaAGWStage
}

//KarnaDeploymentPrune => Karna model for Prune option in Karna config file.
type KarnaDeploymentPrune struct {
	Alias bool `json:"alias"`
	Keep  int  `json:"keep"`
}

//KarnaDeployment => Karna model for Deployment key in Karna config file.
type KarnaDeployment struct {
	Src          string               `json:"src"`
	Key          string               `json:"key"`
	File         string               `json:"file"`
	FunctionName string               `json:"functionName"`
	Aliases      map[string]string    `json:"aliases"`
	Bucket       string               `json:"bucket"`
	Prune        KarnaDeploymentPrune `json:"prune"`
}

//KarnaConfigFile => Karna model for Karna config file.
type KarnaConfigFile struct {
	Global      map[string]string `json:"global"`
	Deployments []KarnaDeployment `json:"deployments"`
	Path        string
}

//KarnaQuery => Karna model for Neo4J query.
type KarnaQuery struct {
	Queries     []string
	QueriesChan chan []string
	Args        []map[string]interface{}
	ArgsChan    chan []map[string]interface{}
}
