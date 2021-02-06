package deploy

import (
	"context"
	"fmt"
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
func (karnaLambdaModel *KarnaLambdaModel) PublishFunction(functionName string, deployment *KarnaDeployment) (err error) {
	input := &lambda.PublishVersionInput{
		FunctionName: aws.String(functionName),
	}

	req := karnaLambdaModel.Client.PublishVersionRequest(input)

	_, err = req.Send(context.Background())

	return
}

//UpdateFunctionCode => Expose UpdateFunctionCode to KarnaLambdaModel.
func (karnaLambdaModel *KarnaLambdaModel) UpdateFunctionCode(deployment *KarnaDeployment, archivePath, functionName string) (err error) {
	var input lambda.UpdateFunctionCodeInput

	if deployment.Bucket != "" {
		input = lambda.UpdateFunctionCodeInput{
			FunctionName: aws.String(functionName),
			S3Bucket:     aws.String(deployment.Bucket),
			S3Key:        aws.String(deployment.File),
			Publish:      aws.Bool(true),
		}
	} else {
		part, _ := ioutil.ReadFile(archivePath)

		input = lambda.UpdateFunctionCodeInput{
			FunctionName: aws.String(functionName),
			Publish:      aws.Bool(true),
			ZipFile:      part,
		}
	}

	req := karnaLambdaModel.Client.UpdateFunctionCodeRequest(&input)

	function, err := req.Send(context.Background())

	if err != nil {
		return
	}

	logger.Log("Current version: " + *function.Version)
	return
}

//GetFunctionByFunctionName => Expose GetFunctionByFunctionName to KarnaLambdaModel.
func (karnaLambdaModel *KarnaLambdaModel) getFunctionByFunctionName(functionName string) (result *lambda.GetFunctionConfigurationResponse, err error) {
	input := &lambda.GetFunctionConfigurationInput{
		FunctionName: aws.String(functionName),
	}

	req := karnaLambdaModel.Client.GetFunctionConfigurationRequest(input)

	result, err = req.Send(context.Background())

	return
}

//GetVersionsByFunction => Expose GetVersionsByFunction to KarnaLambdaModel.
func (karnaLambdaModel *KarnaLambdaModel) getVersionsByFunction(functionName string) (versions []lambda.FunctionConfiguration, err error) {
	input := &lambda.ListVersionsByFunctionInput{
		FunctionName: aws.String(functionName),
	}

	req := karnaLambdaModel.Client.ListVersionsByFunctionRequest(input)

	response, err := req.Send(context.Background())

	if err != nil {
		return
	}

	versions = response.Versions

	return
}

//GetAliasesByFunctionName => Expose GetAliasesByFunctionName to KarnaLambdaModel.
func (karnaLambdaModel *KarnaLambdaModel) getAliasesByFunctionName(functionName string) (aliases []lambda.AliasConfiguration, err error) {
	input := &lambda.ListAliasesInput{
		FunctionName: aws.String(functionName),
	}

	req := karnaLambdaModel.Client.ListAliasesRequest(input)

	response, err := req.Send(context.Background())

	if err != nil {
		return
	}

	aliases = response.Aliases

	return
}

//SyncAlias => Expose SyncAlias to KarnaLambdaModel.
func (karnaLambdaModel *KarnaLambdaModel) syncAlias(deployment *KarnaDeployment, alias, functionName string) (err error) {

	aliases, _ := karnaLambdaModel.getAliasesByFunctionName(functionName)

	if a := findAlias(aliases, alias); a == nil {
		logger.Log("Creating alias: " + alias)
		err = karnaLambdaModel.createAlias(deployment, alias, functionName)
	} else {
		logger.Log("Updating alias: " + alias)
		err = karnaLambdaModel.updateAlias(deployment, alias, functionName)
	}

	return
}

func (karnaLambdaModel *KarnaLambdaModel) createAlias(deployment *KarnaDeployment, alias, functionName string) (err error) {
	var version string

	if len(deployment.Aliases[alias]) == 0 {
		version = "$LATEST"
	} else if deployment.Aliases[alias] == "fixed" {
		versions, _ := karnaLambdaModel.getVersionsByFunction(functionName)
		version = *versions[len(versions)-1].Version
	} else {
		version = deployment.Aliases[alias]
	}

	input := &lambda.CreateAliasInput{
		FunctionName:    aws.String(functionName),
		Name:            aws.String(alias),
		FunctionVersion: aws.String(version),
	}

	req := karnaLambdaModel.Client.CreateAliasRequest(input)

	_, err = req.Send(context.Background())

	return
}

func (karnaLambdaModel *KarnaLambdaModel) updateAlias(deployment *KarnaDeployment, alias, functionName string) (err error) {
	var version string
	functions, err := karnaLambdaModel.getVersionsByFunction(functionName)

	if err != nil {
		return
	}

	if deployment.Aliases[alias] == "fixed" {
		version = *functions[len(functions)-1].Version
	} else {
		if ok := findVersion(functions, deployment.Aliases[alias]); !ok {
			return fmt.Errorf("Version specified do not exists, operation aborted")
		}

		version = deployment.Aliases[alias]
	}

	input := &lambda.UpdateAliasInput{
		FunctionName:    aws.String(functionName),
		Name:            aws.String(alias),
		FunctionVersion: aws.String(version),
	}

	req := karnaLambdaModel.Client.UpdateAliasRequest(input)

	_, err = req.Send(context.Background())

	return
}

func (karnaLambdaModel *KarnaLambdaModel) deleteAlias(functionName, alias string, deployment *KarnaDeployment) (err error) {
	input := &lambda.DeleteAliasInput{
		Name:         aws.String(alias),
		FunctionName: aws.String(functionName),
	}

	req := karnaLambdaModel.Client.DeleteAliasRequest(input)
	_, err = req.Send(context.Background())
	return
}

/**
* -|1|-|2|-|3|-|dev:4|-|5|-|6|-|7|-|prod:8|-|9|-|latest:10|
* => keep 2 - alias prod, versions removed => [1]
* => keep 1 - alias dev, versions removed => [1,2,6]
 */
func (karnaLambdaModel *KarnaLambdaModel) removeVersions(functionName string, deployment *KarnaDeployment) (err error) {
	if deployment.Versions.Keep > 0 {
		var versionsWithAliases []int
		var versionsToPrune []int
		var versionsToKeep []int

		versions, _ := karnaLambdaModel.getVersionsByFunction(functionName)
		aliases, _ := karnaLambdaModel.getAliasesByFunctionName(functionName)

		for _, alias := range aliases {
			version, _ := strconv.Atoi(*alias.FunctionVersion)
			versionsWithAliases = append(versionsWithAliases, version)
		}

		sort.Ints(versionsWithAliases)

		if deployment.Versions.From == "each" {
			for _, v := range versionsWithAliases {
				step := deployment.Versions.Keep
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
		} else {
			minVersion := versionsWithAliases[0]

			for _, f := range versions {
				version, err := strconv.Atoi(*f.Version)

				if err == nil {
					if version < (minVersion - deployment.Versions.Keep) {
						versionsToPrune = append(versionsToPrune, version)
					}
				}
			}
		}

		pruneVersionsCount := strconv.Itoa(len(versionsToPrune))

		logger.Log(pruneVersionsCount + " version(s) removed")

		var wg sync.WaitGroup

		for _, version := range versionsToPrune {
			wg.Add(1)
			karnaLambdaModel.pruneVersion(&wg, version, functionName)
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

func (karnaLambdaModel *KarnaLambdaModel) addPermission(functionURI string) (result *lambda.AddPermissionResponse, err error) {
	input := &lambda.AddPermissionInput{
		FunctionName: aws.String(functionURI),
		Action:       aws.String("lambda:InvokeFunction"),
		Principal:    aws.String("apigateway.amazonaws.com"),
		StatementId:  aws.String(strconv.Itoa(int(time.Now().UnixNano()))),
	}

	req := karnaLambdaModel.Client.AddPermissionRequest(input)

	result, err = req.Send(context.Background())

	return
}
