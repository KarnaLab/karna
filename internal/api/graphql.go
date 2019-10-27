package api

import (
	"encoding/json"
	"github.com/karnalab/karna/core"
	"net/http"

	"github.com/graphql-go/graphql"
)

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})

	if len(result.Errors) > 0 {
		core.LogErrorMessage("An error occured.")
	}

	return result
}

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"lambda": &graphql.Field{
				Type:    graphql.NewList(KarnaGraphQLLambdaType),
				Resolve: lambdaResolver,
			},
			"ec2": &graphql.Field{
				Type:    KarnaGraphQLEC2Type,
				Resolve: ec2Resolver,
			},
			"apigateway": &graphql.Field{
				Type:    graphql.NewList(KarnaGraphQLAGWType),
				Resolve: agwResolver,
			},
		},
	})

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queryType,
	},
)

func buildGraphQLAPI(w http.ResponseWriter, r *http.Request) {
	result := executeQuery(r.URL.Query().Get("query"), schema)

	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
}
