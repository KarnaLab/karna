package viz

import (
	"karna/core"
	"sync"
)

const (
	baseLambdaQuery = "MERGE (lambda:Lambda { uuid: {FunctionArn}, FunctionArn: {FunctionArn}, FunctionName: {FunctionName} }) WITH lambda"

	layersQuery = `
	UNWIND {layers} AS layer
	MERGE (la:Layer { uuid: layer.Arn })
	WITH DISTINCT la, lambda
	MERGE (lambda)-[:HAS_LAYER]->(la)
	WITH DISTINCT lambda
	`
	versionsQuery = `
	UNWIND {versions} AS version
	MERGE (v:Version {uuid: version.uuid, FunctionArn: version.FunctionArn, Version: version.Version, FunctionName: version.FunctionName})
	WITH DISTINCT v, lambda
	MERGE (lambda)-[:HAS_VERSION]->(v)
	WITH DISTINCT lambda
	`
	vpcQuery = `
	MERGE (vpc:VPC {uuid: {vpc}, ID: {vpc}})
	WITH lambda, vpc
	MERGE (vpc)<-[:BELONGS_TO_VPC]-(lambda)
	WITH lambda
	`
	apiGatewayQuery = `
	UNWIND {apiGatewayIds} AS id
	MERGE (apiGateway:ApiGateway {uuid: id })
	WITH DISTINCT apiGateway, lambda
	MERGE (apiGateway)<-[:HAS_TRIGGER]-(lambda)
	WITH DISTINCT lambda
	`
	s3Query = `
	UNWIND {s3Buckets} AS bucket
	MERGE (s3:S3 {uuid: bucket, bucket: bucket })
	WITH DISTINCT s3, lambda
	MERGE (s3)<-[:HAS_TRIGGER]-(lambda)
	WITH DISTINCT lambda
	`
	cloudWatchQuery = `
	UNWIND {cloudWatchRules} AS rule
	MERGE (cloudWatch:CloudWatch {uuid: rule, name: rule })
	WITH DISTINCT cloudWatch, lambda
	MERGE (cloudWatch)<-[:HAS_TRIGGER]-(lambda)
	`
	endQuery = " RETURN lambda"
)

func buildLambdaGraph(wg *sync.WaitGroup) {
	var query = core.Query{
		Args:        []map[string]interface{}{},
		Queries:     []string{},
		ArgsChan:    make(chan []map[string]interface{}),
		QueriesChan: make(chan []string),
	}

	lambdaTree := core.Lambda.BuildLambdaTree()

	go buildLambdaQuery(&query, lambdaTree)

	<-query.QueriesChan
	<-query.ArgsChan

	core.Bulk(query.Queries, query.Args)
	wg.Done()
}

func buildLambdaQuery(query *core.Query, functions []core.KarnaLambda) {
	for _, function := range functions {
		var versions []map[string]interface{}
		var layers []map[string]interface{}

		for _, version := range function.Versions {
			versions = append(versions, map[string]interface{}{
				"uuid":         *version.FunctionArn,
				"FunctionArn":  *version.FunctionArn,
				"Version":      *version.Version,
				"FunctionName": *version.FunctionName,
			})
		}

		for _, layer := range function.Layers {
			layers = append(layers, map[string]interface{}{
				"Arn": *layer.Arn,
			})
		}

		query.Queries = append(query.Queries, buildRequest(function))
		query.Args = append(query.Args, map[string]interface{}{
			"uuid":            *function.FunctionConfiguration.FunctionArn,
			"FunctionArn":     *function.FunctionConfiguration.FunctionArn,
			"FunctionName":    *function.FunctionConfiguration.FunctionName,
			"layers":          layers,
			"vpc":             function.VPC,
			"versions":        versions,
			"apiGatewayIds":   function.Policy["APIGateway"],
			"s3Buckets":       function.Policy["S3"],
			"cloudWatchRules": function.Policy["CloudWatch"],
		})
	}

	query.QueriesChan <- query.Queries
	query.ArgsChan <- query.Args
}

func buildRequest(function core.KarnaLambda) (query string) {
	query = baseLambdaQuery

	if len(function.VPC) > 0 {
		query = query + vpcQuery
	}

	if len(function.Layers) > 0 {
		query = query + layersQuery
	}

	if len(function.Versions) > 0 {
		query = query + versionsQuery
	}

	if len(function.Policy["APIGateway"]) > 0 {
		query = query + apiGatewayQuery
	}

	if len(function.Policy["S3"]) > 0 {
		query = query + s3Query
	}

	if len(function.Policy["CloudWatch"]) > 0 {
		query = query + cloudWatchQuery
	}

	query = query + endQuery

	return
}
