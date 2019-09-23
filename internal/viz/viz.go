package viz

import (
	"karna/core"
	"sync"
	"time"
)

func Run() (timeElapsed string) {
	var wg sync.WaitGroup
	startTime := time.Now()

	wg.Add(3)

	go buildLambdaGraph(&wg)
	go buildAGWGraph(&wg)
	go buildEC2Tree(&wg)

	wg.Wait()

	timeElapsed = time.Since(startTime).String()
	return
}

func Cleanup() (timeElapsed string) {
	startTime := time.Now()
	core.CleanUp()

	timeElapsed = time.Since(startTime).String()
	return
}
