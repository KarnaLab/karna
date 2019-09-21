package viz

import (
	"karna/core"
)

const (
	baseAGWQuery = `
		MERGE (apiGateway:ApiGateway {uuid: { Id }, Name: { Name } }) 
		WITH apiGateway
	`
	startStagesQuery = `
		FOREACH (stage IN {stages} | 
			MERGE (sta:Stage {
				uuid: stage.Uuid, 
				Name: stage.Name
			})<-[:HAS_STAGE]-(apiGateway)
	`
	endStagesQuery = `
		FOREACH (resource IN {resources} |
			MERGE (agwResource:ApiGatewayResource { 
				uuid: resource.uuid, 
				Id: resource.Id,
				Path: resource.Path
			})
			MERGE (sta)-[:HAS_RESOURCE]->(agwResource)
		)
	)	
	WITH apiGateway
	`
	cloudFrontDistributionStageQuery = `
			MERGE (distribution:CloudFrontDistribution { 
				uuid: stage.Distribution, 
				Name: stage.Distribution 
			})
			MERGE (sta)-[:HAS_DISTRIBUTION]->(distribution)
	`
	//TEMP hack
	cleanupQuery = `
		MATCH (cfd:CloudFrontDistribution) WHERE cfd.uuid=''
		DETACH DELETE cfd
	`
)

//TODO: Need to get domain names && modelize them into Neo4j.
func buildAGWQuery(query *core.Query, apis []core.KarnaAGWAPI) {
	for _, api := range apis {
		var hasDistribution bool
		var stages []map[string]interface{}
		request := baseAGWQuery + startStagesQuery

		for _, stage := range api.Stages {
			if len(stage.Distribution) > 0 {
				hasDistribution = true
			}
		}
		if hasDistribution {
			request = request + cloudFrontDistributionStageQuery
		}
		for _, stage := range api.Stages {
			stages = append(stages, map[string]interface{}{
				"Name":         stage.Name,
				"Uuid":         stage.Uuid,
				"Distribution": stage.Distribution,
			})
		}

		request = request + endStagesQuery + cleanupQuery
		query.Queries = append(query.Queries, request)
		query.Args = append(query.Args, map[string]interface{}{
			"Id":        *api.API.Id,
			"Name":      *api.API.Name,
			"resources": api.Resources,
			"stages":    stages,
			//"distribution": api.Distribution,
		})
	}

	query.QueriesChan <- query.Queries
	query.ArgsChan <- query.Args
}

func buildAGWGraph() {
	var query = core.Query{
		Args:        []map[string]interface{}{},
		Queries:     []string{},
		ArgsChan:    make(chan []map[string]interface{}),
		QueriesChan: make(chan []string),
	}
	AGWTree := core.AGW.BuildAGWTree()

	go buildAGWQuery(&query, AGWTree)

	<-query.QueriesChan
	<-query.ArgsChan

	core.Bulk(query.Queries, query.Args)
}
