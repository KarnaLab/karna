package deploy

import (
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
)

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
	UUID         string
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
	Alias bool `json:"alias,omitempty"`
	Keep  int  `json:"keep,omitempty"`
}

type KarnaAPIDeployment struct {
	ID         string `json:"id,omitempty"`
	Resource   string `json:"resource,omitempty"`
	HTTPMethod string `json:"httpMethod,omitempty"`
}

//KarnaDeployment => Karna model for Deployment key in Karna config file.
type KarnaDeployment struct {
	Src          string               `json:"src"`
	Key          string               `json:"key,omitempty"`
	File         string               `json:"file"`
	FunctionName string               `json:"functionName"`
	Aliases      map[string]string    `json:"aliases,omitempty"`
	Bucket       string               `json:"bucket,omitempty"`
	Prune        KarnaDeploymentPrune `json:"prune,omitempty"`
	Executable   string               `json:"executable,omitempty"`
	API          KarnaAPIDeployment   `json:"api,omitempty"`
}

//KarnaConfigFile => Karna model for Karna config file.
type KarnaConfigFile struct {
	Global      map[string]string `json:"global"`
	Deployments []KarnaDeployment `json:"deployments"`
	Path        string            `json:",omitempty"`
}
