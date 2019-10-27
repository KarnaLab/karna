package deploy

import (
	"testing"

	"github.com/karnalab/karna/core"
)

func TestCheckRequirements(t *testing.T) {
	alias := "dev"
	deploymentTest := core.KarnaDeployment{
		Src:          "src",
		Key:          "key",
		File:         "file",
		FunctionName: "functionName",
		Aliases: map[string]string{
			"dev": "version",
		},
	}
	checkRequirements(&deploymentTest, alias)
}

func TestGetConfigFile(t *testing.T) {
	configFile := getConfigFile()

	if configFile != nil {
		t.Error("biaarre")
	}
}
