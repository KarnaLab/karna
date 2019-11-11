package core

import (
	"errors"
	"os"
	"reflect"

	"github.com/spf13/viper"
)

const (
	fileName = "karna"
)

func GetConfig() (config *KarnaConfigFile, err error) {

	viper.SetConfigName(fileName)
	viper.AddConfigPath(".")

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	viper.Unmarshal(&config)
	// Check if all required keys are provided:
	err = checkRequirements(config)

	if err != nil {
		return
	}

	// Merge global config into functions config:
	err = mergeGlobalConfig(config)

	if err != nil {
		return
	}

	return
}

func checkRequirements(config *KarnaConfigFile) (err error) {
	requirements := [...]string{"Name", "Aliases", "Input"}

	for _, function := range config.Functions {
		for _, requirement := range requirements {
			a := reflect.Indirect(reflect.ValueOf(function)).FieldByName(requirement)
			switch a.Type() {
			case reflect.TypeOf(""):
				if a.Len() == 0 {
					err = errors.New("is missing:" + requirement)
					return
				}
			case reflect.TypeOf(map[string]string{}):
				if a.IsNil() {
					err = errors.New("is missing:" + requirement)
					return
				}
			}
		}
	}

	return
}

func mergeGlobalConfig(config *KarnaConfigFile) (err error) {
	dir, err := os.Getwd()
	hasGlobalOutputDefined := viper.IsSet("global.output")

	for i, function := range config.Functions {
		hasFunctionOutput := len(function.Output) > 0

		if hasFunctionOutput {
			continue
		}

		if !hasFunctionOutput && !hasGlobalOutputDefined {
			return errors.New("Output is missing either in Global config && Function config")
		}
		config.Functions[i].Output = config.Global.Output
	}

	config.Path = dir

	return
}
