package viz

import (
	"karna/core"
	"sync"
	"time"
)

//Run => Will build all AWS dependencies into trees and load them into Neo4J.
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

//Cleanup => Will detach delete all Neo4J nodes.
func Cleanup() (timeElapsed string) {
	startTime := time.Now()
	core.CleanUp()

	timeElapsed = time.Since(startTime).String()
	return
}
