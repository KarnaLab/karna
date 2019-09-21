package viz

import (
	"awsgraphviz/utils"
	"karna/core"
	"sync"
)

const (
	request = `
        MERGE (vpc:VPC {uuid: { VpcId }, ID: { VpcId } })
        WITH vpc
        UNWIND {Subnets} AS subnet
        MERGE (sub:Subnet {uuid: subnet.uuid, Name: subnet.Name, CidrBlock: subnet.CidrBlock })
        MERGE (sub)-[:BELONGS_TO_VPC]->(vpc)
        WITH DISTINCT subnet, vpc
        MERGE (instance:EC2Instance { uuid: { uuid }, Name: { Name } })
        WITH instance
        MATCH (s:Subnet {uuid: { SubnetId } })
        MERGE (s)<-[:HAS_SUBNET]-(instance)
        WITH instance
        MATCH (v:VPC {uuid: { VpcId } })
        MERGE (v)<-[:BELONGS_TO_VPC]-(instance)
    `
)

func buildEC2Query(query *core.Query, ec2 core.KarnaEC2Model) {
	var subnets []map[string]interface{}

	for _, subnet := range ec2.Subnets {
		subnets = append(subnets, map[string]interface{}{
			"uuid":      *subnet.SubnetId,
			"Name":      *subnet.SubnetId,
			"VpcId":     *subnet.VpcId,
			"CidrBlock": *subnet.CidrBlock,
		})
	}

	for _, instance := range ec2.Instances {
		query.Queries = append(query.Queries, request)
		query.Args = append(query.Args, map[string]interface{}{
			"uuid":     *instance.InstanceId,
			"Name":     *instance.InstanceId,
			"SubnetId": *instance.SubnetId,
			"VpcId":    *instance.VpcId,
			"Subnets":  subnets,
		})
	}

	query.QueriesChan <- query.Queries
	query.ArgsChan <- query.Args
}

func buildEC2Tree(wg *sync.WaitGroup) {
	var query = core.Query{
		Args:        []map[string]interface{}{},
		Queries:     []string{},
		ArgsChan:    make(chan []map[string]interface{}),
		QueriesChan: make(chan []string),
	}

	ec2Tree := core.EC2.BuildEC2Tree()

	go buildEC2Query(&query, ec2Tree)

	<-query.QueriesChan
	<-query.ArgsChan

	utils.Bulk(query.Queries, query.Args)
	wg.Done()
}
