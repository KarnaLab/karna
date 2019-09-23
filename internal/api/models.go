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

var KarnaGraphQLAGWType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "APIGateway",
		Fields: graphql.Fields{
			"API": &graphql.Field{
				Type: KarnaGraphQLAGWRestAPIType,
			},
			"Resources": &graphql.Field{
				Type: graphql.String,
			},
			"Stages": &graphql.Field{
				Type: graphql.NewList(KarnaGraphQLAGWStageType),
			},
		},
	},
)

/*
    ApiKeySource    ApiKeySourceType    `locationName:"apiKeySource" type:"string" enum:"true"`
    BinaryMediaTypes    []string    `locationName:"binaryMediaTypes" type:"list"`
    CreatedDate    *time.Time    `locationName:"createdDate" type:"timestamp"`
    Description    *string    `locationName:"description" type:"string"`
    EndpointConfiguration    *EndpointConfiguration    `locationName:"endpointConfiguration" type:"structure"`
    Id    *string    `locationName:"id" type:"string"`
    MinimumCompressionSize    *int64    `locationName:"minimumCompressionSize" type:"integer"`
    Name    *string    `locationName:"name" type:"string"`
    Policy    *string    `locationName:"policy" type:"string"`
    Tags    map[string]string    `locationName:"tags" type:"map"`
    Version    *string    `locationName:"version" type:"string"`
		Warnings    []string    `locationName:"warnings" type:"list"`
*/
var KarnaGraphQLAGWRestAPIType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RestAPI",
		Fields: graphql.Fields{
			"ApiKeySource": &graphql.Field{
				Type: graphql.String,
			},
			"BinaryMediaTypes": &graphql.Field{
				Type: graphql.String,
			},
			"CreatedDate": &graphql.Field{
				Type: graphql.String,
			},
			"Description": &graphql.Field{
				Type: graphql.String,
			},
			"Id": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
var KarnaGraphQLAGWStageType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Layer",
		Fields: graphql.Fields{
			"Name": &graphql.Field{
				Type: graphql.String,
			},
			"Stage": &graphql.Field{
				Type: graphql.String,
			},
			"FunctionName": &graphql.Field{
				Type: graphql.String,
			},
			"UUID": &graphql.Field{
				Type: graphql.String,
			},
			"Distribution": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
