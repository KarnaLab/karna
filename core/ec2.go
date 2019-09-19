package core

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func (KarnaEC2 *KarnaEC2) init() {
	cfg, err := external.LoadDefaultAWSConfig()

	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	KarnaEC2.Client = ec2.New(cfg)
}

func (KarnaEC2 *KarnaEC2) BuildEC2Tree() KarnaEC2Model {

	instances := KarnaEC2.getInstances()
	activeSubnets, securityGroups := KarnaEC2.fetchDependencies()
	vpcs := KarnaEC2.getVPCS(instances)

	modelizedEC2 := KarnaEC2Model{
		Instances:      instances,
		SecurityGroups: securityGroups,
		Subnets:        activeSubnets,
		VPCS:           vpcs,
	}

	return modelizedEC2
}

func (KarnaEC2 *KarnaEC2) fetchDependencies() (activeSubnets []ec2.Subnet, securityGroups []ec2.SecurityGroup) {
	activeSubnetsChan := make(chan []ec2.Subnet, 1)
	securityGroupsChan := make(chan []ec2.SecurityGroup, 1)

	go KarnaEC2.getActiveSubnets(activeSubnetsChan)
	go KarnaEC2.getSecurityGroups(securityGroupsChan)

	activeSubnets = <-activeSubnetsChan
	securityGroups = <-securityGroupsChan

	return
}

func (KarnaEC2 *KarnaEC2) getInstances() (instances []ec2.Instance) {
	input := &ec2.DescribeInstancesInput{}

	req := KarnaEC2.Client.DescribeInstancesRequest(input)

	results, err := req.Send(context.Background())

	if err != nil {
		panic(err.Error())
	}

	for _, reservation := range results.Reservations {
		for _, instance := range reservation.Instances {
			instances = append(instances, instance)
		}
	}

	return
}

func (KarnaEC2 *KarnaEC2) getActiveSubnets(activeSubnetsChan chan []ec2.Subnet) {
	var activeSubnets []ec2.Subnet
	input := &ec2.DescribeSubnetsInput{}
	req := KarnaEC2.Client.DescribeSubnetsRequest(input)
	results, _ := req.Send(context.Background())

	for _, subnet := range results.Subnets {
		activeSubnets = append(activeSubnets, subnet)
	}
	activeSubnetsChan <- activeSubnets
}

func (KarnaEC2 *KarnaEC2) getSecurityGroups(securityGroupsChan chan []ec2.SecurityGroup) {
	var securityGroups []ec2.SecurityGroup
	input := &ec2.DescribeSecurityGroupsInput{}
	req := KarnaEC2.Client.DescribeSecurityGroupsRequest(input)
	results, _ := req.Send(context.Background())

	for _, securityGroup := range results.SecurityGroups {
		securityGroups = append(securityGroups, securityGroup)
	}

	securityGroupsChan <- securityGroups
}

func (KarnaEC2 *KarnaEC2) getVPCS(instances []ec2.Instance) (VPCS []string) {
	var vpcs []string

	for _, instance := range instances {
		vpcs = append(vpcs, *instance.VpcId)
	}

	VPCS = uniq(vpcs)

	return
}
