package viz

import (
	"strings"
	"sync"
	"time"

	"github.com/karnalab/karna/core"
)

var neo4j core.KarnaNeo4J

func parseCredentials(credentials string) (username, password string) {

	parsed := strings.Split(credentials, "/")

	if len(parsed) > 1 {
		username = parsed[0]
		password = parsed[1]
	}

	return
}

//Run => Will build all AWS dependencies into trees and load them into Neo4J.
func Run(port, credentials, host *string) (timeElapsed string) {
	var wg sync.WaitGroup
	startTime := time.Now()
	username, password := parseCredentials(*credentials)

	neo4j.Configuration = core.KarnaNeo4JConfiguration{
		Username: username,
		Password: password,
		Port:     *port,
		Host:     *host,
	}

	wg.Add(3)

	go buildLambdaGraph(&wg)
	go buildAGWGraph(&wg)
	go buildEC2Tree(&wg)

	wg.Wait()

	timeElapsed = time.Since(startTime).String()
	return
}

//Cleanup => Will detach delete all Neo4J nodes.
func Cleanup(port, credentials, host *string) (timeElapsed string) {
	startTime := time.Now()
	username, password := parseCredentials(*credentials)

	neo4j.Configuration = core.KarnaNeo4JConfiguration{
		Username: username,
		Password: password,
		Port:     *port,
		Host:     *host,
	}

	neo4j.CleanUp()

	timeElapsed = time.Since(startTime).String()
	return
}
