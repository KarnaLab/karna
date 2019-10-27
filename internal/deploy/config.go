package deploy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/karnalab/karna/core"
)

const (
	fileName = "karna.json"
)

func getConfigFile() (configFile *core.KarnaConfigFile, err error) {
	dir, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(dir + "/" + fileName)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &configFile)

	configFile.Path = dir

	if err != nil {
		return nil, err
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

func checkRequirements(deployment *core.KarnaDeployment, alias string) (err error) {
	requirements := [...]string{"FunctionName", "File", "Aliases", "Src"}

	if deployment.Aliases[alias] == "" {
		return fmt.Errorf("alias do not match with the config file")
	}

	for _, requirement := range requirements {

		a := reflect.Indirect(reflect.ValueOf(deployment)).FieldByName(requirement)

		switch a.Type() {
		case reflect.TypeOf(""):
			if a.Len() == 0 {
				return fmt.Errorf("is missing:" + requirement)
			}
		case reflect.TypeOf(map[string]string{}):
			if a.IsNil() {
				return fmt.Errorf("is missing:" + requirement)
			}
		}
	}
	return
}
