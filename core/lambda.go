package core

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go/aws"
)

//BuildLambdaTree => Will Build Lambda tree for Karna model.
func (lambdaModel *KarnaLambdas) BuildLambdaTree() []KarnaLambda {
	var wg sync.WaitGroup
	functions := lambdaModel.getFunctions()

	modelizedFunctions := make([]KarnaLambda, len(functions))

	for i, function := range functions {
		var vpc string
		wg.Add(1)

		if function.VpcConfig != nil && len(*function.VpcConfig.VpcId) > 0 {
			vpc = *function.VpcConfig.VpcId
		}

		modelizedFunctions[i] = KarnaLambda{
			FunctionConfiguration: function,
			Layers:                function.Layers,
			VPC:                   vpc,
		}
		go lambdaModel.fetchDependencies(&modelizedFunctions[i], &wg)
	}

	wg.Wait()

	return modelizedFunctions
}

func (lambdaModel *KarnaLambdas) init() {
	cfg, err := external.LoadDefaultAWSConfig()

	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	lambdaModel.Client = lambda.New(cfg)
}

func (lambdaModel *KarnaLambdas) fetchDependencies(function *KarnaLambda, wg *sync.WaitGroup) {
	versions := make(chan []lambda.FunctionConfiguration, 1)
	policy := make(chan map[string][]string, 1)

	go lambdaModel.getVersions(versions, *function.FunctionConfiguration.FunctionName)
	go lambdaModel.getPolicy(policy, *function.FunctionConfiguration.FunctionArn)

	function.Versions = <-versions
	function.Policy = <-policy

	wg.Done()
}

func (lambdaModel *KarnaLambdas) getFunctions() (functions []lambda.FunctionConfiguration) {
	input := &lambda.ListFunctionsInput{}

	req := lambdaModel.Client.ListFunctionsRequest(input)

	response, err := req.Send(context.Background())

	if err != nil {
		panic(err.Error())
	}

	functions = response.Functions

	return
}

func (lambdaModel *KarnaLambdas) getVersions(versions chan []lambda.FunctionConfiguration, functionName string) {
	var listVersionsInput interface{}

	listVersionsInput = &lambda.ListVersionsByFunctionInput{FunctionName: aws.String(functionName)}
	request := lambdaModel.Client.ListVersionsByFunctionRequest(listVersionsInput.(*lambda.ListVersionsByFunctionInput))

	response, err := request.Send(context.Background())

	if err != nil {
		panic(err.Error())
	}

	versions <- response.Versions
}

func (lambdaModel *KarnaLambdas) getPolicy(policies chan map[string][]string, functionArn string) {
	var policyInput interface{}
	var policy awsPolicy
	dependencies := make(map[string][]string, 1)

	policyInput = &lambda.GetPolicyInput{FunctionName: aws.String(functionArn)}
	request := lambdaModel.Client.GetPolicyRequest(policyInput.(*lambda.GetPolicyInput))

	response, _ := request.Send(context.Background())

	if response != nil {
		json.Unmarshal([]byte(*response.Policy), &policy)

		for _, statement := range policy.Statement {
			switch statement.Principal.Service {
			case "apigateway.amazonaws.com":
				dependencies["APIGateway"] = append(dependencies["APIGateway"], findAPIGatewayID(statement))
			case "s3.amazonaws.com":
				dependencies["S3"] = append(dependencies["S3"], findS3Bucket(statement))
			case "events.amazonaws.com":
				dependencies["CloudWatch"] = append(dependencies["CloudWatch"], findCloudWatch(statement))
			default:
				fmt.Println("Unhandled service")
			}
		}
	}

	policies <- dependencies
}
