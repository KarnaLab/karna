package core

import (
	"context"
	"karna/core"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go/aws"
)

func (agw *KarnaAPIGateway) init() {
	cfg, err := external.LoadDefaultAWSConfig()

	if err != nil {
		core.LogErrorMessage("unable to load SDK config, " + err.Error())
		os.Exit(2)
	}

	agw.Client = apigateway.New(cfg)
}

//BuildAGWTree => Will Build APIGateway tree for Karna model.
func (agw *KarnaAPIGateway) BuildAGWTree() []KarnaAGWAPI {
	var wg sync.WaitGroup
	apis := agw.getAPIS()

	modelizedAPIS := make([]KarnaAGWAPI, len(apis))

	for i, api := range apis {
		wg.Add(1)

		modelizedAPIS[i] = KarnaAGWAPI{
			API: api,
		}
		go agw.fetchDependencies(&modelizedAPIS[i], &wg)
	}

	wg.Wait()

	for i := range modelizedAPIS {
		wg.Add(1)
		go agw.fetchPathMappings(&modelizedAPIS[i], &modelizedAPIS, &wg)
	}

	wg.Wait()

	return modelizedAPIS
}

func (agw *KarnaAPIGateway) fetchDependencies(api *KarnaAGWAPI, wg *sync.WaitGroup) {
	resources := make(chan []map[string]interface{}, 1)
	stages := make(chan []KarnaAGWStage, 1)

	go agw.getResourcesForAPI(resources, *api.API.Id)
	go agw.getStagesByAPI(stages, *api.API.Id)

	api.Resources = <-resources
	api.Stages = <-stages

	wg.Done()
}

func (agw *KarnaAPIGateway) getAPIS() (apis []apigateway.RestApi) {
	input := &apigateway.GetRestApisInput{}

	req := agw.Client.GetRestApisRequest(input)

	results, err := req.Send(context.Background())

	if err != nil {
		core.LogErrorMessage(err.Error())
		os.Exit(2)
	}

	apis = results.Items

	return
}

func (agw *KarnaAPIGateway) getStagesByAPI(stagesChan chan []KarnaAGWStage, id string) {
	var stages []KarnaAGWStage
	input := &apigateway.GetStagesInput{RestApiId: aws.String(id)}
	req := agw.Client.GetStagesRequest(input)

	results, err := req.Send(context.Background())

	if err != nil {
		core.LogErrorMessage(err.Error())
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

func (agw *KarnaAPIGateway) getResourcesForAPI(resourcesChan chan []map[string]interface{}, id string) {
	var resources []map[string]interface{}
	input := &apigateway.GetResourcesInput{RestApiId: aws.String(id)}
	req := agw.Client.GetResourcesRequest(input)

	results, err := req.Send(context.Background())

	if err != nil {
		core.LogErrorMessage(err.Error())
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

func (agw *KarnaAPIGateway) fetchPathMappings(api *KarnaAGWAPI, apis *[]KarnaAGWAPI, wg *sync.WaitGroup) {
	input := &apigateway.GetDomainNamesInput{}
	req := agw.Client.GetDomainNamesRequest(input)

	results, _ := req.Send(context.Background())

	for _, domainName := range results.Items {
		mappings := agw.getBasePathMappings(domainName)

		for _, mapping := range mappings {
			if *mapping.RestApiId == *api.API.Id {
				stageIndex := findStage(api.Stages, *mapping.Stage)
				api.Stages[stageIndex].Distribution = *domainName.DistributionDomainName
			}
		}
	}

	wg.Done()
}

func (agw *KarnaAPIGateway) getBasePathMappings(domainName apigateway.DomainName) (mappings []apigateway.BasePathMapping) {
	input := &apigateway.GetBasePathMappingsInput{DomainName: aws.String(*domainName.DomainName)}
	req := agw.Client.GetBasePathMappingsRequest(input)

	results, _ := req.Send(context.Background())
	mappings = results.Items

	return
}
