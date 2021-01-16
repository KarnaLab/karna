package deploy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
)

const (
	fileName = "karna.json"
)

func getConfigFile() (configFile *KarnaConfigFile, err error) {
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

func getTargetDeployment(config *KarnaConfigFile, target *string) (deployment *KarnaDeployment, err error) {
	for _, d := range config.Deployments {
		if d.FunctionName == *target {
			deployment = &d
		}
	}

	if deployment == nil {
		err = fmt.Errorf("Deployment not found in config file")
	}

	return
}

func checkRequirements(deployment *KarnaDeployment, alias string) (err error) {
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
