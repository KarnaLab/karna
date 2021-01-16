package deploy

import (
	"context"
	"io/ioutil"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go/aws"
)

type KarnaLambdaModel struct {
	Client *lambda.Client
}

func (karnaLambdaModel *KarnaLambdaModel) init() {
	cfg, err := external.LoadDefaultAWSConfig()

	if err != nil {
		logger.Error("unable to load SDK config, " + err.Error())
	}

	karnaLambdaModel.Client = lambda.New(cfg)
}

//PublishFunction => Expose PublishFunction to KarnaLambdaModel.
func (karnaLambdaModel *KarnaLambdaModel) PublishFunction(deployment *KarnaDeployment) (err error) {
	input := &lambda.PublishVersionInput{
		FunctionName: aws.String(deployment.FunctionName),
	}

	req := karnaLambdaModel.Client.PublishVersionRequest(input)

	_, err = req.Send(context.Background())

	return
}

//UpdateFunctionCode => Expose UpdateFunctionCode to KarnaLambdaModel.
func (karnaLambdaModel *KarnaLambdaModel) UpdateFunctionCode(deployment *KarnaDeployment, archivePath string) (err error) {
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

	req := karnaLambdaModel.Client.UpdateFunctionCodeRequest(input)

	_, err = req.Send(context.Background())

	return
}

//GetFunctionByFunctionName => Expose GetFunctionByFunctionName to KarnaLambdaModel.
func (karnaLambdaModel *KarnaLambdaModel) GetFunctionByFunctionName(functionName string) (result *lambda.GetFunctionConfigurationResponse, err error) {
	input := &lambda.GetFunctionConfigurationInput{
		FunctionName: aws.String(functionName),
	}

	req := karnaLambdaModel.Client.GetFunctionConfigurationRequest(input)

	result, err = req.Send(context.Background())

	return
}

//GetVersionsByFunction => Expose GetVersionsByFunction to KarnaLambdaModel.
func (karnaLambdaModel *KarnaLambdaModel) GetVersionsByFunction(functionName string) (versions []lambda.FunctionConfiguration, err error) {
	input := &lambda.ListVersionsByFunctionInput{
		FunctionName: aws.String(functionName),
	}

	req := karnaLambdaModel.Client.ListVersionsByFunctionRequest(input)

	response, err := req.Send(context.Background())

	versions = response.Versions

	return
}

//GetAliasesByFunctionName => Expose GetAliasesByFunctionName to KarnaLambdaModel.
func (karnaLambdaModel *KarnaLambdaModel) GetAliasesByFunctionName(functionName string) (aliases []lambda.AliasConfiguration, err error) {
	input := &lambda.ListAliasesInput{
		FunctionName: aws.String(functionName),
	}

	req := karnaLambdaModel.Client.ListAliasesRequest(input)

	response, err := req.Send(context.Background())

	aliases = response.Aliases

	return
}

//SyncAlias => Expose SyncAlias to KarnaLambdaModel.
func (karnaLambdaModel *KarnaLambdaModel) SyncAlias(deployment *KarnaDeployment, alias string) (err error) {

	aliases, _ := karnaLambdaModel.GetAliasesByFunctionName(deployment.FunctionName)

	if a := findAlias(aliases, alias); a == nil {
		logger.Log("creation of alias: " + alias)
		karnaLambdaModel.createAlias(deployment, alias)
	} else {
		logger.Log("updating alias: " + alias)
		karnaLambdaModel.updateAlias(deployment, alias)
	}

	return
}

func (karnaLambdaModel *KarnaLambdaModel) createAlias(deployment *KarnaDeployment, alias string) (err error) {
	var version string

	if len(deployment.Aliases[alias]) == 0 {
		version = "$LATEST"
	} else if deployment.Aliases[alias] == "fixed" {
		versions, _ := karnaLambdaModel.GetVersionsByFunction(deployment.FunctionName)
		version = *versions[len(versions)-1].Version
	} else {
		version = deployment.Aliases[alias]
	}

	input := &lambda.CreateAliasInput{
		FunctionName:    aws.String(deployment.FunctionName),
		Name:            aws.String(alias),
		FunctionVersion: aws.String(version),
	}

	req := karnaLambdaModel.Client.CreateAliasRequest(input)

	_, err = req.Send(context.Background())

	return
}

func (karnaLambdaModel *KarnaLambdaModel) updateAlias(deployment *KarnaDeployment, alias string) (err error) {
	var version string

	if deployment.Aliases[alias] == "fixed" {
		versions, _ := karnaLambdaModel.GetVersionsByFunction(deployment.FunctionName)
		version = *versions[len(versions)-1].Version
	} else {
		version = deployment.Aliases[alias]
	}

	input := &lambda.UpdateAliasInput{
		FunctionName:    aws.String(deployment.FunctionName),
		Name:            aws.String(alias),
		FunctionVersion: aws.String(version),
	}

	req := karnaLambdaModel.Client.UpdateAliasRequest(input)

	_, err = req.Send(context.Background())

	return
}

//Prune => Expose Prune to KarnaLambdaModel. Will remove alias and/or versions.
func (karnaLambdaModel *KarnaLambdaModel) Prune(deployment *KarnaDeployment) (err error) {
	if deployment.Prune.Alias {
		aliases, _ := karnaLambdaModel.GetAliasesByFunctionName(deployment.FunctionName)

		for _, a := range aliases {
			if _, ok := deployment.Aliases[*a.Name]; !ok {
				logger.Log("prune alias: " + *a.Name)

				input := &lambda.DeleteAliasInput{
					Name:         aws.String(*a.Name),
					FunctionName: aws.String(deployment.FunctionName),
				}

				req := karnaLambdaModel.Client.DeleteAliasRequest(input)
				_, err = req.Send(context.Background())
			}
		}
	}

	if deployment.Prune.Keep > 0 {
		var versionsWithAliases []int
		var versionsToPrune []int
		var versionsToKeep []int

		versions, _ := karnaLambdaModel.GetVersionsByFunction(deployment.FunctionName)
		aliases, _ := karnaLambdaModel.GetAliasesByFunctionName(deployment.FunctionName)

		for _, alias := range aliases {
			version, _ := strconv.Atoi(*alias.FunctionVersion)
			versionsWithAliases = append(versionsWithAliases, version)
		}

		sort.Ints(versionsWithAliases)

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

		pruneVersionsCount := strconv.Itoa(len(versionsToPrune))

		logger.Log("prune: " + pruneVersionsCount + " version(s)")

		var wg sync.WaitGroup

		for _, version := range versionsToPrune {
			wg.Add(1)
			karnaLambdaModel.pruneVersion(&wg, version, deployment.FunctionName)
		}

		wg.Wait()
	}

	return
}

func (karnaLambdaModel *KarnaLambdaModel) pruneVersion(wg *sync.WaitGroup, version int, functionName string) {
	versionToString := strconv.Itoa(version)

	input := &lambda.DeleteFunctionInput{
		FunctionName: aws.String(functionName),
		Qualifier:    aws.String(versionToString),
	}

	req := karnaLambdaModel.Client.DeleteFunctionRequest(input)

	_, err := req.Send(context.Background())

	if err != nil {
		logger.Error(err.Error())
	}

	wg.Done()
}

func (karnaLambdaModel *KarnaLambdaModel) AddPermission(functionName, alias string) (result *lambda.AddPermissionResponse, err error) {
	input := &lambda.AddPermissionInput{
		FunctionName: aws.String(functionName + ":" + alias),
		Action:       aws.String("lambda:InvokeFunction"),
		Principal:    aws.String("apigateway.amazonaws.com"),
		StatementId:  aws.String(strconv.Itoa(int(time.Now().UnixNano()))),
	}

	req := karnaLambdaModel.Client.AddPermissionRequest(input)

	result, err = req.Send(context.Background())

	return
}
