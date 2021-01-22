package deploy

import (
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/viper"
)

const (
	fileName = "karna"
)

func getConfigFile() (configFile *KarnaConfigFile, err error) {
	viper.SetConfigName(fileName)
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	dir, err := os.Getwd()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	viper.Unmarshal(&configFile)

	if err != nil {
		return nil, err
	}

	configFile.Path = dir

	return
}

func getTargetDeployment(config *KarnaConfigFile, functionName string) (deployment *KarnaDeployment, err error) {
	if d, ok := config.Deployments[functionName]; ok {
		deployment = &d
	} else {
		err = fmt.Errorf("Deployment not found in config file")
	}
	return
}

func checkRequirements(deployment *KarnaDeployment, alias string) (err error) {
	requirements := [...]string{"File", "Aliases", "Src"}

	if deployment.Aliases[alias] == "" {
		return fmt.Errorf("Alias do not match with the config file")
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
