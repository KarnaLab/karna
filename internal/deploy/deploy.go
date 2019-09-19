package deploy

import (
	"fmt"
	"karna/core"
)

func Run(target *string, alias *string) {

	fmt.Println(*target, *alias)

	configFile := getConfigFile()
	targetDeployment := getTargetDeployment(configFile)

	checkRequirements(targetDeployment)

	zipArchive(configFile.Path+"/"+targetDeployment.Src, configFile.Path+"/.karna/"+targetDeployment.File)

	core.S3.Upload(configFile.Path + "/.karna/" + targetDeployment.File)
}
