package deploy

import (
	"encoding/json"
	"io/ioutil"
	"karna/core"
	"os"
	"reflect"
)

const (
	fileName = "karna.json"
)

func getConfigFile() (configFile *core.KarnaConfigFile) {
	dir, err := os.Getwd()

	if err != nil {
		core.LogErrorMessage(err.Error())
		os.Exit(2)
	}

	data, err := ioutil.ReadFile(dir + "/" + fileName)

	if err != nil {
		core.LogErrorMessage(err.Error())
		os.Exit(2)
	}

	err = json.Unmarshal(data, &configFile)

	configFile.Path = dir

	if err != nil {
		core.LogErrorMessage(err.Error())
		os.Exit(2)
	}

	return
}

func getTargetDeployment(config *core.KarnaConfigFile, target *string) (deployment *core.KarnaDeployment) {
	for _, d := range config.Deployments {
		if d.FunctionName == *target {
			deployment = &d
		}
	}
	return
}

func checkRequirements(deployment *core.KarnaDeployment, alias string) {
	requirements := [...]string{"FunctionName", "File", "Aliases", "Src"}

	if deployment.Aliases[alias] == "" {
		core.LogErrorMessage("Alias do not match with the config file.")
		os.Exit(2)
	}

	for _, requirement := range requirements {

		a := reflect.Indirect(reflect.ValueOf(deployment)).FieldByName(requirement)

		switch a.Type() {
		case reflect.TypeOf(""):
			if a.Len() == 0 {
				core.LogErrorMessage("is missing:" + requirement)
				os.Exit(2)
			}
		case reflect.TypeOf(map[string]string{}):
			if a.IsNil() {
				core.LogErrorMessage("is missing:" + requirement)
				os.Exit(2)
			}
		}
	}
}
