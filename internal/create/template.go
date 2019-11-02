package create

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/karnalab/karna/core"
)

const (
	indexTemplate = `exports.handler = async (event) => {
    const response = {
        statusCode: 200,
        body: JSON.stringify('Hello from Lambda!'),
    };
    return response;
};
	`
	filePerm   = 0644
	folderPerm = 0755
)

func createFolder(folder string) (err error) {
	err = os.Mkdir(folder, folderPerm)
	return
}

func createFileWithTemplate(file, template string) (err error) {
	data := []byte(template)
	err = ioutil.WriteFile(file, data, filePerm)

	return
}

func generateDeploymentConfig(functionName *string) (deployment *core.KarnaDeployment) {
	deployment = &core.KarnaDeployment{
		Src:          *functionName,
		File:         "lambda.zip",
		FunctionName: *functionName,
		Aliases: map[string]string{
			"dev":  "fixed@update",
			"prod": "1",
		},
	}
	return
}

func generateKarnaConfigFile(folder, functionName *string) (err error) {
	config := &core.KarnaConfigFile{
		Global:      map[string]string{},
		Deployments: []core.KarnaDeployment{},
	}

	path := *folder + "/karna.json"
	deployment := generateDeploymentConfig(functionName)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		config.Deployments = append(config.Deployments, *deployment)
	} else {
		data, err := ioutil.ReadFile(path)

		if err != nil {
			return err
		}

		json.Unmarshal(data, &config)

		isDeploymentDefined := core.FindDeployment(*functionName, config.Deployments)

		if !isDeploymentDefined {
			config.Deployments = append(config.Deployments, *deployment)

			os.Remove(path)
		}
	}

	jsonData, err := json.Marshal(config)

	if err != nil {
		return err
	}

	createFileWithTemplate(path, string(jsonData))
	return
}

func generateLayout(folder, functionName *string) {

	if _, err := os.Stat(*folder); os.IsNotExist(err) {
		createFolder(*folder)
	}

	generateKarnaConfigFile(folder, functionName)

	functionFolder := *folder + "/" + *functionName

	if _, err := os.Stat(functionFolder); os.IsNotExist(err) {
		createFolder(functionFolder)
	}

	if _, err := os.Stat(functionFolder + "/index.js"); os.IsNotExist(err) {
		createFileWithTemplate(functionFolder+"/index.js", indexTemplate)
	}
}

func generateLayersTemplate(runtime, folder *string) {

	if _, err := os.Stat(*folder + "/common"); os.IsNotExist(err) {
		createFolder(*folder + "/common")
	}

	switch *runtime {
	case "nodejs":
		generateNodeJSRuntime(folder)
		break
	}
}

func generateNodeJSRuntime(folder *string) {
	path := *folder + "/common/nodejs"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		createFolder(path)
	}

	if _, err := os.Stat(path + "/package.json"); os.IsNotExist(err) {
		createFileWithTemplate(path+"/package.json", "{}")
	}
}
