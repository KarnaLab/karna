package deploy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"karna/core"
	"log"
	"os"
	"reflect"
)

const (
	fileName = "karna.json"
)

func getConfigFile() (configFile *core.KarnaConfigFile) {
	dir, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadFile(dir + "/" + fileName)

	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(data, &configFile)

	configFile.Path = dir

	if err != nil {
		fmt.Println(err)
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

func checkRequirements(deployment *core.KarnaDeployment) {
	requirements := [...]string{"FunctionName", "Key", "File", "Aliases", "Bucket"}
	for _, requirement := range requirements {

		a := reflect.Indirect(reflect.ValueOf(deployment)).FieldByName(requirement)

		switch a.Type() {
		case reflect.TypeOf(""):
			if a.Len() == 0 {
				panic("is missing:" + requirement)
			}
		case reflect.TypeOf(map[string]string{}):
			if a.IsNil() {
				panic("is missing:" + requirement)
			}
		}
	}
}
