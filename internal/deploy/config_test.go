package deploy

import (
	"testing"

	"github.com/karnalab/karna/core"
	"github.com/karnalab/karna/create"
)

func TestCheckRequirementsWithValidInputs(t *testing.T) {
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
	err := checkRequirements(&deploymentTest, alias)

	if err == nil {
		t.Log("Test PASSED with the right alias")
	} else {
		t.Errorf(err.Error())
	}
}

func TestCheckRequirementsWithInvalidInputs(t *testing.T) {
	alias := "wrongAlias"
	deploymentTest := core.KarnaDeployment{
		Src:          "src",
		Key:          "key",
		File:         "file",
		FunctionName: "functionName",
		Aliases: map[string]string{
			"dev": "version",
		},
	}
	err := checkRequirements(&deploymentTest, alias)

	if err != nil {
		t.Log("Test FAILED because the wrong alias")
	} else {
		t.Errorf("checkRequirement must not find alias in KarnaDeployment")
	}
}
func TestGetConfigFileWithoutConfigFile(t *testing.T) {
	_, err := getConfigFile()

	if err != nil {
		t.Log(err.Error())
	} else {
		t.Errorf("getConfigFile must not find config file")
	}
}

func TestGetConfigFileWithConfigFile(t *testing.T) {
	create.Run("test", "functionName", "nodejs")
	_, err := getConfigFile()

	if err != nil {
		t.Log(err.Error())
	} else {
		t.Errorf("getConfigFile must not find config file")
	}
}
