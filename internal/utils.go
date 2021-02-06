package deploy

import (
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

func findInt(needle int, haystack []int) (found bool) {
	for _, value := range haystack {
		if needle == value {
			found = true
		}
	}
	return
}

func findAlias(aliases []lambda.AliasConfiguration, aliasName string) (alias *lambda.AliasConfiguration) {
	for _, a := range aliases {
		if *a.Name == aliasName {
			alias = &a
		}
	}
	return
}

func findVersion(functions []lambda.FunctionConfiguration, version string) (found bool) {
	for _, a := range functions {
		if *a.Version == version {
			found = true
			return
		}
	}
	return
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}
