package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
)

func executeQuery(query string, schema graphql.Schema) (result *graphql.Result, err error) {
	result = graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})

	if len(result.Errors) > 0 {
		return result, fmt.Errorf("an error occured")
	}

	return
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
	result, err := executeQuery(r.URL.Query().Get("query"), schema)

	if err != nil {
		panic(err)
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
}
