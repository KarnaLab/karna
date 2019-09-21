package viz

import (
	"karna/core"
	"sync"
)

func Run() {
	var wg sync.WaitGroup

	go buildLambdaGraph(&wg)
	go buildAGWGraph(&wg)
	go buildEC2Tree(&wg)

	wg.Wait()
}

func Cleanup() {
	core.CleanUp()
}
