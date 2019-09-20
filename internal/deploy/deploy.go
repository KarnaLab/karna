package deploy

import (
	"fmt"
	"karna/core"
)

func Run(target *string, alias *string) {

	configFile := getConfigFile()
	targetDeployment := getTargetDeployment(configFile, target)

	checkRequirements(targetDeployment, *alias)

	var source = configFile.Path + "/" + targetDeployment.Src
	var output = configFile.Path + "/.karna/" + targetDeployment.File

	zipArchive(source, output)

	if targetDeployment.Bucket != "" {
		err := core.S3.Upload(targetDeployment, output)

		if err != nil {
			panic(err.Error())
		}
	}

	err := core.Lambda.UpdateFunctionCode(targetDeployment, output)

	if err != nil {
		fmt.Println(err.Error())
	}

	err = core.Lambda.SyncAlias(targetDeployment, *alias)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err.Error())
	}

	//Prune functions

}
