package core

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func (karnaEC2Model *KarnaEC2Model) init() {
	cfg, err := external.LoadDefaultAWSConfig()

	if err != nil {
		logger.Error("unable to load SDK config, " + err.Error())
	}

	karnaEC2Model.Client = ec2.New(cfg)
}

//BuildEC2Tree => will build a tree with all EC2, activeSubnets,securityGroups && VPCS associated.
func (karnaEC2Model *KarnaEC2Model) BuildEC2Tree() KarnaEC2 {

	instances := karnaEC2Model.getInstances()
	activeSubnets, securityGroups := karnaEC2Model.fetchDependencies()
	vpcs := getVPCS(instances)

	modelizedEC2 := KarnaEC2{
		Instances:      instances,
		SecurityGroups: securityGroups,
		Subnets:        activeSubnets,
		VPCS:           vpcs,
	}

	return modelizedEC2
}

func (karnaEC2Model *KarnaEC2Model) fetchDependencies() (activeSubnets []ec2.Subnet, securityGroups []ec2.SecurityGroup) {
	activeSubnetsChan := make(chan []ec2.Subnet, 1)
	securityGroupsChan := make(chan []ec2.SecurityGroup, 1)

	go karnaEC2Model.getActiveSubnets(activeSubnetsChan)
	go karnaEC2Model.getSecurityGroups(securityGroupsChan)

	activeSubnets = <-activeSubnetsChan
	securityGroups = <-securityGroupsChan

	return
}

func (karnaEC2Model *KarnaEC2Model) getInstances() (instances []ec2.Instance) {
	input := &ec2.DescribeInstancesInput{}

	req := karnaEC2Model.Client.DescribeInstancesRequest(input)

	results, err := req.Send(context.Background())

	if err != nil {
		logger.Error(err.Error())
	}

	for _, reservation := range results.Reservations {
		for _, instance := range reservation.Instances {
			instances = append(instances, instance)
		}
	}

	return
}

func (karnaEC2Model *KarnaEC2Model) getActiveSubnets(activeSubnetsChan chan []ec2.Subnet) {
	var activeSubnets []ec2.Subnet

	input := &ec2.DescribeSubnetsInput{}
	req := karnaEC2Model.Client.DescribeSubnetsRequest(input)
	results, err := req.Send(context.Background())

	if err != nil {
		logger.Error(err.Error())
	}

	for _, subnet := range results.Subnets {
		activeSubnets = append(activeSubnets, subnet)
	}
	activeSubnetsChan <- activeSubnets
}

func (karnaEC2Model *KarnaEC2Model) getSecurityGroups(securityGroupsChan chan []ec2.SecurityGroup) {
	var securityGroups []ec2.SecurityGroup

	input := &ec2.DescribeSecurityGroupsInput{}
	req := karnaEC2Model.Client.DescribeSecurityGroupsRequest(input)
	results, err := req.Send(context.Background())

	if err != nil {
		logger.Error(err.Error())
	}

	for _, securityGroup := range results.SecurityGroups {
		securityGroups = append(securityGroups, securityGroup)
	}

	securityGroupsChan <- securityGroups
}

func getVPCS(instances []ec2.Instance) (VPCS []string) {
	var vpcs []string

	for _, instance := range instances {
		vpcs = append(vpcs, *instance.VpcId)
	}

	VPCS = uniq(vpcs)

	return
}
