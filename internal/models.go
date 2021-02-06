package deploy

import (
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
)

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

//KarnaDeploymentVersions => Karna model for Prune option in Karna config file.
type KarnaDeploymentVersions struct {
	Keep int    `json:"keep,omitempty"`
	From string `json:"from,omitempty"`
}

type KarnaAPIDeployment struct {
	ID         string `json:"id,omitempty"`
	Resource   string `json:"resource,omitempty"`
	HTTPMethod string `json:"httpMethod,omitempty"`
}

//KarnaDeployment => Karna model for Deployment key in Karna config file.
type KarnaDeployment struct {
	Src        string                  `json:"src"`
	Key        string                  `json:"key,omitempty"`
	File       string                  `json:"file"`
	Aliases    map[string]string       `json:"aliases,omitempty"`
	Bucket     string                  `json:"bucket,omitempty"`
	Versions   KarnaDeploymentVersions `json:"versions,omitempty"`
	Executable string                  `json:"executable,omitempty"`
	API        KarnaAPIDeployment      `json:"api,omitempty"`
}

//KarnaConfigFile => Karna model for Karna config file.
type KarnaConfigFile struct {
	Global      map[string]string          `json:"global"`
	Deployments map[string]KarnaDeployment `json:"deployments"`
	Path        string                     `json:",omitempty"`
}
