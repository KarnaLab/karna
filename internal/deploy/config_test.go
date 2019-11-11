package deploy

import (
	"testing"

	"github.com/karnalab/karna/core"
)

const (
	filePerm = 0644
)

func TestGetTargetFunctionWithValidInputs(t *testing.T) {
	alias := "dev"
	target := "functionName"

	function := core.KarnaFunction{
		Input:  "src",
		Output: "lambda.zip",
		Name:   "functionName",
		Aliases: map[string]string{
			"dev": "version",
		},
	}

	config := core.KarnaConfigFile{
		Functions: []core.KarnaFunction{
			function,
		},
	}
	_, err := getTargetFunction(&config, &target, &alias)

	if err == nil {
		t.Log("Test PASSED with the right alias and functionName")
	} else {
		t.Errorf(err.Error())
	}
}

func TestGetTargetFunctionWithInvalidInputs(t *testing.T) {
	alias := "invalid"
	target := "functionName"

	function := core.KarnaFunction{
		Input:  "src",
		Output: "lambda.zip",
		Name:   "functionName",
		Aliases: map[string]string{
			"dev": "version",
		},
	}

	config := core.KarnaConfigFile{
		Functions: []core.KarnaFunction{
			function,
		},
	}
	_, err := getTargetFunction(&config, &target, &alias)

	if err != nil {
		t.Log("Test FAILED with the wrong alias")
	} else {
		t.Errorf(err.Error())
	}

	alias = "dev"
	target = "invalid"

	_, err = getTargetFunction(&config, &target, &alias)

	if err != nil {
		t.Log("Test FAILED with the wrong functionName")
	} else {
		t.Errorf(err.Error())
	}
}
