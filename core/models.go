package core

import "github.com/aws/aws-sdk-go-v2/service/lambda"

type KarnaLambdas struct {
	Client *lambda.Client
}

type KarnaLambda struct {
	FunctionConfiguration lambda.FunctionConfiguration
	Layers                []lambda.Layer
	VPC                   string
	Versions              []lambda.FunctionConfiguration
	Policy                map[string][]string
}
