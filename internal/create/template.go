package create

import (
	"karna/core"
	"os"
)

func createFolder(folder string) {
	err := os.Mkdir(folder, 0755)

	if err != nil {
		core.LogErrorMessage(err.Error())
	}
}

func generateTemplate(name, layerName, runtime *string, withLayers *bool) {
	dir, err := os.Getwd()

	if err != nil {
		core.LogErrorMessage(err.Error())
	}

	folder := dir + "/" + *name

	createFolder(folder)

	if *withLayers {
		generateLayersTemplate(runtime, &folder)
	}
}

func generateLayersTemplate(runtime, folder *string) {
	createFolder(*folder + "/layers")

	switch *runtime {
	case "nodejs":
		createFolder(*folder + "/layers/nodejs")
		break
	}
}
