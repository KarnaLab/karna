package create

import (
	"io/ioutil"
	"karna/core"
	"os"
)

const (
	karnaTemplate = "{\"global\":{}, \"deployments\":[]}"
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

func createFolder(folder string) {
	err := os.Mkdir(folder, folderPerm)

	if err != nil {
		core.LogErrorMessage(err.Error())
	}
}

func createFileWithTemplate(file, template string) {
	data := []byte(template)
	err := ioutil.WriteFile(file, data, filePerm)

	if err != nil {
		core.LogErrorMessage(err.Error())
	}
}

func generateTemplate(name, functionName, runtime *string) {
	dir, err := os.Getwd()

	if err != nil {
		core.LogErrorMessage(err.Error())
	}

	folder := dir + "/" + *name

	generateLayout(&folder, functionName)

	if len(*runtime) > 0 {
		generateLayersTemplate(runtime, &folder)
	}
}

func generateKarnaConfigFile(folder *string) {
	path := *folder + "/karna.json"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		createFileWithTemplate(path, karnaTemplate)
	}
}

func generateLayout(folder, functionName *string) {

	if _, err := os.Stat(*folder); os.IsNotExist(err) {
		createFolder(*folder)
	}

	generateKarnaConfigFile(folder)

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
}
