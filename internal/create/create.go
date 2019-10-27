package create

import (
	"github.com/karnalab/karna/core"
	"os"
	"time"
)

func Run(name, functionName, runtime *string) (timeElapsed string) {
	startTime := time.Now()
	dir, err := os.Getwd()

	if err != nil {
		core.LogErrorMessage(err.Error())
	}

	folder := dir + "/" + *name

	generateLayout(&folder, functionName)

	if len(*runtime) > 0 {
		generateLayersTemplate(runtime, &folder)
	}

	timeElapsed = time.Since(startTime).String()
	return
}
