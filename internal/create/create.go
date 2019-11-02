package create

import (
	"os"
	"time"
)

func Run(name, functionName, runtime *string) (timeElapsed string, err error) {
	startTime := time.Now()
	dir, err := os.Getwd()

	if err != nil {
		return timeElapsed, err
	}

	folder := dir + "/" + *name

	generateLayout(&folder, functionName)

	if len(*runtime) > 0 {
		generateLayersTemplate(runtime, &folder)
	}

	timeElapsed = time.Since(startTime).String()
	return
}
