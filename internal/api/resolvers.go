package api

import (
	"github.com/karnalab/karna/core"

	"github.com/graphql-go/graphql"
)

func lambdaResolver(p graphql.ResolveParams) (interface{}, error) {
	response := core.Lambda.BuildLambdaTree()

	return response, nil
}

func ec2Resolver(p graphql.ResolveParams) (interface{}, error) {
	response := core.EC2.BuildEC2Tree()

	return response, nil
}

func agwResolver(p graphql.ResolveParams) (interface{}, error) {
	response := core.AGW.BuildAGWTree()

	return response, nil
}
