package viz

import (
	"karna/core"
)

func Run() {
	buildLambdaGraph()
}

func Cleanup() {
	core.CleanUp()
}
