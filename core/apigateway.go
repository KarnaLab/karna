package core

import (
	"context"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go/aws"
)

func (karnaAGW *KarnaAPIGatewayModel) init() {
	cfg, err := external.LoadDefaultAWSConfig()

	if err != nil {
		LogErrorMessage("unable to load SDK config, " + err.Error())
		os.Exit(2)
	}

	karnaAGW.Client = apigateway.New(cfg)
}

//BuildAGWTree => Will Build APIGateway tree for Karna model.
func (karnaAGW *KarnaAPIGatewayModel) BuildAGWTree() []KarnaAGWAPI {
	var wg sync.WaitGroup
	apis := karnaAGW.getAPIS()

	modelizedAPIS := make([]KarnaAGWAPI, len(apis))

	for i, api := range apis {
		wg.Add(1)

		modelizedAPIS[i] = KarnaAGWAPI{
			API: api,
		}
		go karnaAGW.fetchDependencies(&modelizedAPIS[i], &wg)
	}

	wg.Wait()

	for i := range modelizedAPIS {
		wg.Add(1)
		go karnaAGW.fetchPathMappings(&modelizedAPIS[i], &modelizedAPIS, &wg)
	}

	wg.Wait()

	return modelizedAPIS
}

func (karnaAGW *KarnaAPIGatewayModel) fetchDependencies(api *KarnaAGWAPI, wg *sync.WaitGroup) {
	resources := make(chan []map[string]interface{}, 1)
	stages := make(chan []KarnaAGWStage, 1)

	go karnaAGW.getResourcesForAPI(resources, *api.API.Id)
	go karnaAGW.getStagesByAPI(stages, *api.API.Id)

	api.Resources = <-resources
	api.Stages = <-stages

	wg.Done()
}

func (karnaAGW *KarnaAPIGatewayModel) getAPIS() (apis []apigateway.RestApi) {
	input := &apigateway.GetRestApisInput{}

	req := karnaAGW.Client.GetRestApisRequest(input)

	results, err := req.Send(context.Background())

	if err != nil {
		LogErrorMessage(err.Error())
		os.Exit(2)
	}

	apis = results.Items

	return
}

func (karnaAGW *KarnaAPIGatewayModel) getStagesByAPI(stagesChan chan []KarnaAGWStage, id string) {
	var stages []KarnaAGWStage
	input := &apigateway.GetStagesInput{RestApiId: aws.String(id)}
	req := karnaAGW.Client.GetStagesRequest(input)

	results, err := req.Send(context.Background())

	if err != nil {
		LogErrorMessage(err.Error())
		os.Exit(2)
	}

	for _, stage := range results.Item {
		stages = append(stages, KarnaAGWStage{
			Name:  *stage.StageName,
			Uuid:  id + *stage.StageName,
			Stage: *stage.StageName,
		})
	}

	stagesChan <- stages
}

func (karnaAGW *KarnaAPIGatewayModel) getResourcesForAPI(resourcesChan chan []map[string]interface{}, id string) {
	var resources []map[string]interface{}
	input := &apigateway.GetResourcesInput{RestApiId: aws.String(id)}
	req := karnaAGW.Client.GetResourcesRequest(input)

	results, err := req.Send(context.Background())

	if err != nil {
		LogErrorMessage(err.Error())
		os.Exit(2)
	}

	for _, resource := range results.Items {
		resources = append(resources, map[string]interface{}{
			"Id":   *resource.Id,
			"Path": *resource.Path,
			"uuid": *resource.Id,
		})
	}
	resourcesChan <- resources
}

func (karnaAGW *KarnaAPIGatewayModel) fetchPathMappings(api *KarnaAGWAPI, apis *[]KarnaAGWAPI, wg *sync.WaitGroup) {
	input := &apigateway.GetDomainNamesInput{}
	req := karnaAGW.Client.GetDomainNamesRequest(input)

	results, _ := req.Send(context.Background())

	for _, domainName := range results.Items {
		mappings := karnaAGW.getBasePathMappings(domainName)

		for _, mapping := range mappings {
			if *mapping.RestApiId == *api.API.Id {
				stageIndex := findStage(api.Stages, *mapping.Stage)
				api.Stages[stageIndex].Distribution = *domainName.DistributionDomainName
			}
		}
	}

	wg.Done()
}

func (karnaAGW *KarnaAPIGatewayModel) getBasePathMappings(domainName apigateway.DomainName) (mappings []apigateway.BasePathMapping) {
	input := &apigateway.GetBasePathMappingsInput{DomainName: aws.String(*domainName.DomainName)}
	req := karnaAGW.Client.GetBasePathMappingsRequest(input)

	results, _ := req.Send(context.Background())
	mappings = results.Items

	return
}
