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
