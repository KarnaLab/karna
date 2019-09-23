package api

import "github.com/graphql-go/graphql"

var KarnaGraphQLLayersType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Layer",
		Fields: graphql.Fields{
			"Arn": &graphql.Field{
				Type: graphql.String,
			},
			"CodeSize": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)
var KarnaGraphQLFunctionConfigurationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Layer",
		Fields: graphql.Fields{
			"CodeSha256": &graphql.Field{
				Type: graphql.String,
			},
			"CodeSize": &graphql.Field{
				Type: graphql.Int,
			},
			"FunctionArn": &graphql.Field{
				Type: graphql.String,
			},
			"FunctionName": &graphql.Field{
				Type: graphql.String,
			},
			"Handler": &graphql.Field{
				Type: graphql.String,
			},
			"LastModified": &graphql.Field{
				Type: graphql.String,
			},
			"RevisionId": &graphql.Field{
				Type: graphql.String,
			},
			"MemorySize": &graphql.Field{
				Type: graphql.Int,
			},
			"Role": &graphql.Field{
				Type: graphql.String,
			},
			"Runtime": &graphql.Field{
				Type: graphql.String,
			},
			"Version": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var KarnaGraphQLPolicyType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Policy",
		Fields: graphql.Fields{
			"APIGateway": &graphql.Field{
				Type: graphql.String,
			},
			"S3": &graphql.Field{
				Type: graphql.String,
			},
			"CloudWatch": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var KarnaGraphQLLambdaType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Lambda",
		Fields: graphql.Fields{
			"FunctionConfiguration": &graphql.Field{
				Type: KarnaGraphQLFunctionConfigurationType,
			},
			"Layers": &graphql.Field{
				Type: graphql.NewList(KarnaGraphQLLayersType),
			},
			"VPC": &graphql.Field{
				Type: graphql.String,
			},
			"Versions": &graphql.Field{
				Type: graphql.NewList(KarnaGraphQLFunctionConfigurationType),
			},
			"Policy": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
