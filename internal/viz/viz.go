package viz

import (
	"karna/core"
)

func Run() {
	buildLambdaGraph()
	buildAGWGraph()
	buildEC2Tree()
}

func Cleanup() {
	core.CleanUp()
}
