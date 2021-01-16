package deploy

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

const (
	filePerm = 0644
)

func TestCheckRequirementsWithValidInputs(t *testing.T) {
	alias := "dev"
	deploymentTest := KarnaDeployment{
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
	deploymentTest := KarnaDeployment{
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
		t.Errorf("getConfigFile FAILED because it must not find config file")
	}
}

func TestGetConfigFileWithConfigFile(t *testing.T) {
	functionName := "test"
	deployment := &KarnaDeployment{
		Src:          functionName,
		File:         "lambda.zip",
		FunctionName: functionName,
		Aliases: map[string]string{
			"dev":  "fixed",
			"prod": "1",
		},
	}
	jsonData, err := json.Marshal(deployment)
	data := []byte(jsonData)
	err = ioutil.WriteFile("./karna.json", data, filePerm)

	_, err = getConfigFile()

	if err == nil {
		t.Log("getConfigFile PASSED because it must find config file")
	} else {
		t.Errorf(err.Error())
	}
	os.Remove("./karna.json")
}

func TestGetTargetDeploymentWithCorrectTarget(t *testing.T) {
	deployment := KarnaDeployment{
		Src:          "functionName",
		File:         "lambda.zip",
		FunctionName: "functionName",
		Aliases: map[string]string{
			"dev":  "fixed",
			"prod": "1",
		},
	}
	config := KarnaConfigFile{
		Global:      map[string]string{},
		Deployments: []KarnaDeployment{deployment},
	}
	target := "functionName"
	targetDeployment, _ := getTargetDeployment(&config, &target)

	if len(targetDeployment.Src) > 0 {
		t.Log("getTargetDeployment PASSED because it must find the deployment with the right target")
	} else {
		t.Errorf("getTargetDeployment FAILED because it must find the deployment")
	}
}
