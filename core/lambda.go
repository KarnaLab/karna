package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
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

func (lambdaModel *KarnaLambdas) UpdateFunctionCode(deployment *KarnaDeployment, archivePath string) (err error) {

	input := &lambda.UpdateFunctionCodeInput{}

	if deployment.Bucket != "" {
		input = &lambda.UpdateFunctionCodeInput{
			FunctionName: aws.String(deployment.FunctionName),
			S3Bucket:     aws.String(deployment.Bucket),
			S3Key:        aws.String(deployment.File),
			Publish:      aws.Bool(true),
		}
	} else {
		part, _ := ioutil.ReadFile(archivePath)

		input = &lambda.UpdateFunctionCodeInput{
			FunctionName: aws.String(deployment.FunctionName),
			Publish:      aws.Bool(true),
			ZipFile:      part,
		}
	}

	req := lambdaModel.Client.UpdateFunctionCodeRequest(input)

	_, err = req.Send(context.Background())

	return
}

func (lambdaModel *KarnaLambdas) PublishFunction(deployment *KarnaDeployment) (err error) {

	input := &lambda.PublishVersionInput{
		FunctionName: aws.String(deployment.FunctionName),
		Description:  aws.String("gloup"),
	}

	req := lambdaModel.Client.PublishVersionRequest(input)

	_, err = req.Send(context.Background())

	return
}

func (lambdaModel *KarnaLambdas) GetFunctionByFunctionName(functionName string) (err error) {
	input := &lambda.GetFunctionConfigurationInput{
		FunctionName: aws.String(functionName),
	}

	req := lambdaModel.Client.GetFunctionConfigurationRequest(input)

	response, err := req.Send(context.Background())

	fmt.Println(response)
	return
}

func (lambdaModel *KarnaLambdas) GetVersionsByFunction(functionName string) (versions []lambda.FunctionConfiguration, err error) {
	input := &lambda.ListVersionsByFunctionInput{
		FunctionName: aws.String(functionName),
	}

	req := lambdaModel.Client.ListVersionsByFunctionRequest(input)

	response, err := req.Send(context.Background())

	versions = response.Versions

	return
}

func (lambdaModel *KarnaLambdas) GetAliasesByFunctionName(functionName string) (aliases []lambda.AliasConfiguration, err error) {
	input := &lambda.ListAliasesInput{
		FunctionName: aws.String(functionName),
	}

	req := lambdaModel.Client.ListAliasesRequest(input)

	response, err := req.Send(context.Background())

	aliases = response.Aliases

	return
}

func (lambdaModel *KarnaLambdas) SyncAlias(deployment *KarnaDeployment, alias string) (err error) {

	aliases, _ := lambdaModel.GetAliasesByFunctionName(deployment.FunctionName)

	if a := findAlias(aliases, alias); a == nil {
		fmt.Println("create alias")
		lambdaModel.createAlias(deployment, alias)
	} else {
		fmt.Println("update alias")
		lambdaModel.updateAlias(deployment, alias)
	}

	return
}

func (lambdaModel *KarnaLambdas) createAlias(deployment *KarnaDeployment, alias string) (err error) {
	var version string

	if deployment.Aliases[alias] == "fixed@update" || len(deployment.Aliases[alias]) == 0 {
		version = "$LATEST"
	} else {
		version = deployment.Aliases[alias]
	}

	fmt.Println(deployment.Aliases[alias] == "fixed@update", version)
	input := &lambda.CreateAliasInput{
		FunctionName:    aws.String(deployment.FunctionName),
		Name:            aws.String(alias),
		FunctionVersion: aws.String(version),
	}

	req := lambdaModel.Client.CreateAliasRequest(input)

	_, err = req.Send(context.Background())

	return
}

func (lambdaModel *KarnaLambdas) updateAlias(deployment *KarnaDeployment, alias string) (err error) {
	var version string

	if deployment.Aliases[alias] == "fixed@update" {
		versions, _ := lambdaModel.GetVersionsByFunction(deployment.FunctionName)
		version = *versions[len(versions)-1].Version
	} else {
		version = deployment.Aliases[alias]
	}

	input := &lambda.UpdateAliasInput{
		FunctionName:    aws.String(deployment.FunctionName),
		Name:            aws.String(alias),
		FunctionVersion: aws.String(version),
	}

	req := lambdaModel.Client.UpdateAliasRequest(input)

	_, err = req.Send(context.Background())

	return
}

func (lambdaModel *KarnaLambdas) Prune(deployment *KarnaDeployment) (err error) {

	if deployment.Prune.Alias {
		aliases, _ := lambdaModel.GetAliasesByFunctionName(deployment.FunctionName)

		for _, a := range aliases {
			if _, ok := deployment.Aliases[*a.Name]; !ok {
				input := &lambda.DeleteAliasInput{
					Name:         aws.String(*a.Name),
					FunctionName: aws.String(deployment.FunctionName),
				}

				req := lambdaModel.Client.DeleteAliasRequest(input)
				_, err = req.Send(context.Background())
			}
		}
	}

	if deployment.Prune.Keep > 0 {
		var versionsWithAliases []int
		var versionsToPrune []int

		versions, _ := lambdaModel.GetVersionsByFunction(deployment.FunctionName)
		aliases, _ := lambdaModel.GetAliasesByFunctionName(deployment.FunctionName)

		for _, alias := range aliases {
			version, _ := strconv.Atoi(*alias.FunctionVersion)
			versionsWithAliases = append(versionsWithAliases, version)
		}

		sort.Ints(versionsWithAliases)
		var versionsToKeep []int

		for _, v := range versionsWithAliases {
			step := deployment.Prune.Keep
			min := v - step
			max := v + step

			if min <= 1 {
				min = 1
			}
			rangeOfVersions := makeRange(min, max)
			versionsToKeep = append(versionsToKeep, rangeOfVersions...)
		}

		for _, f := range versions {
			version, err := strconv.Atoi(*f.Version)

			if err == nil {
				if ok := findInt(version, versionsToKeep); !ok {
					versionsToPrune = append(versionsToPrune, version)
				}
			}
		}

		for _, version := range versionsToPrune {
			versionToString := strconv.Itoa(version)

			input := &lambda.DeleteFunctionInput{
				FunctionName: aws.String(deployment.FunctionName),
				Qualifier:    aws.String(versionToString),
			}

			req := lambdaModel.Client.DeleteFunctionRequest(input)

			_, err := req.Send(context.Background())

			fmt.Println(err)
		}
	}

	return
}
