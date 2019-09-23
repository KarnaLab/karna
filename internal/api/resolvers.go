package api

import (
	"karna/core"

	"github.com/graphql-go/graphql"
)

var Fields graphql.Fields

func LambdaResolver(p graphql.ResolveParams) (interface{}, error) {
	response := core.Lambda.BuildLambdaTree()

	return response, nil
}

func EC2Resolver(p graphql.ResolveParams) (interface{}, error) {
	response := core.EC2.BuildEC2Tree()

	return response, nil
}

func AGWResolver(p graphql.ResolveParams) (interface{}, error) {
	response := core.AGW.BuildAGWTree()

	return response, nil
}
