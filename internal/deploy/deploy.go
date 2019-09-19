package deploy

import (
	"fmt"
	"karna/core"
)

func Run(target *string, alias *string) {

	fmt.Println(*target, *alias)

	configFile := getConfigFile()
	targetDeployment := getTargetDeployment(configFile, target)

	checkRequirements(targetDeployment)

	var source = configFile.Path + "/" + targetDeployment.Src
	var output = configFile.Path + "/.karna/" + targetDeployment.File

	//TODO: Remove file if exists
	zipArchive(source, output)

	err := core.S3.Upload(targetDeployment, output)

	if err != nil {
		panic(err.Error())
	}

	err = core.Lambda.UpdateFunctionCode(targetDeployment, output)

	if err != nil {
		fmt.Println(err.Error())
	}
	//PublishFunction
	//ProcessAlias
	// - updateAlias or createAlias
	//Prune functions

}
