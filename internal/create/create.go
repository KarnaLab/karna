package create

import (
	"time"

	"github.com/karnalab/karna/core"
)

func Run(functionName, APIName, APIEndpoint, resource, verb *string) (timeElapsed string, err error) {
	var logger *core.KarnaLogger
	startTime := time.Now()

	logger.Log("Checking requirements...")

	if err := checkRequirements(APIEndpoint, verb); err != nil {
		return timeElapsed, err
	}

	timeElapsed = time.Since(startTime).String()
	return
}
