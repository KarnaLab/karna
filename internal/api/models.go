package api

import "github.com/graphql-go/graphql"

var KarnaGraphQLLayersType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Layer",
		Fields: graphql.Fields{
			"Arn": &graphql.Field{
				Type: graphql.String,
			},
			"CodeSize": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)
var KarnaGraphQLFunctionConfigurationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Layer",
		Fields: graphql.Fields{
			"CodeSha256": &graphql.Field{
				Type: graphql.String,
			},
			"CodeSize": &graphql.Field{
				Type: graphql.Int,
			},
			"FunctionArn": &graphql.Field{
				Type: graphql.String,
			},
			"FunctionName": &graphql.Field{
				Type: graphql.String,
			},
			"Handler": &graphql.Field{
				Type: graphql.String,
			},
			"LastModified": &graphql.Field{
				Type: graphql.String,
			},
			"RevisionId": &graphql.Field{
				Type: graphql.String,
			},
			"MemorySize": &graphql.Field{
				Type: graphql.Int,
			},
			"Role": &graphql.Field{
				Type: graphql.String,
			},
			"Runtime": &graphql.Field{
				Type: graphql.String,
			},
			"Version": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var KarnaGraphQLPolicyType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Policy",
		Fields: graphql.Fields{
			"APIGateway": &graphql.Field{
				Type: graphql.String,
			},
			"S3": &graphql.Field{
				Type: graphql.String,
			},
			"CloudWatch": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var KarnaGraphQLLambdaType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Lambda",
		Fields: graphql.Fields{
			"FunctionConfiguration": &graphql.Field{
				Type: KarnaGraphQLFunctionConfigurationType,
			},
			"Layers": &graphql.Field{
				Type: graphql.NewList(KarnaGraphQLLayersType),
			},
			"VPC": &graphql.Field{
				Type: graphql.String,
			},
			"Versions": &graphql.Field{
				Type: graphql.NewList(KarnaGraphQLFunctionConfigurationType),
			},
			"Policy": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var KarnaGraphQLAGWType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "APIGateway",
		Fields: graphql.Fields{
			"API": &graphql.Field{
				Type: KarnaGraphQLAGWRestAPIType,
			},
			"Resources": &graphql.Field{
				Type: graphql.NewList(KarnaGraphQLAGWResourceType),
			},
			"Stages": &graphql.Field{
				Type: graphql.NewList(KarnaGraphQLAGWStageType),
			},
		},
	},
)

var KarnaGraphQLAGWRestAPIType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RestAPI",
		Fields: graphql.Fields{
			"ApiKeySource": &graphql.Field{
				Type: graphql.String,
			},
			"BinaryMediaTypes": &graphql.Field{
				Type: graphql.String,
			},
			"CreatedDate": &graphql.Field{
				Type: graphql.String,
			},
			"Description": &graphql.Field{
				Type: graphql.String,
			},
			"Id": &graphql.Field{
				Type: graphql.String,
			},
			"MinimumCompressionSize": &graphql.Field{
				Type: graphql.Int,
			},
			"Name": &graphql.Field{
				Type: graphql.String,
			},
			"Policy": &graphql.Field{
				Type: graphql.String,
			},
			"Version": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var KarnaGraphQLAGWStageType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Stage",
		Fields: graphql.Fields{
			"Name": &graphql.Field{
				Type: graphql.String,
			},
			"Stage": &graphql.Field{
				Type: graphql.String,
			},
			"FunctionName": &graphql.Field{
				Type: graphql.String,
			},
			"UUID": &graphql.Field{
				Type: graphql.String,
			},
			"Distribution": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var KarnaGraphQLAGWResourceType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Resource",
		Fields: graphql.Fields{
			"Id": &graphql.Field{
				Type: graphql.String,
			},
			"Path": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var KarnaGraphQLEC2Type = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "EC2",
		Fields: graphql.Fields{
			"Instances": &graphql.Field{
				Type: graphql.NewList(KarnaGraphQLInstanceType),
			},
			"SecurityGroups": &graphql.Field{
				Type: graphql.NewList(KarnaGraphQLSecurityGroupType),
			},
			"Subnets": &graphql.Field{
				Type: graphql.NewList(KarnaGraphQLSubnetType),
			},
			"VPCS": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
		},
	},
)

var KarnaGraphQLSubnetType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Subnet",
		Fields: graphql.Fields{
			"AssignIpv6AddressOnCreation": &graphql.Field{
				Type: graphql.Boolean,
			},
			"AvailabilityZone": &graphql.Field{
				Type: graphql.String,
			},
			"AvailabilityZoneId": &graphql.Field{
				Type: graphql.String,
			},
			"AvailableIpAddressCount": &graphql.Field{
				Type: graphql.Int,
			},
			"CidrBlock": &graphql.Field{
				Type: graphql.String,
			},
			"DefaultForAz": &graphql.Field{
				Type: graphql.Boolean,
			},
			"MapPublicIpOnLaunch": &graphql.Field{
				Type: graphql.Boolean,
			},
			"OwnerId": &graphql.Field{
				Type: graphql.Boolean,
			},
			"State": &graphql.Field{
				Type: graphql.String,
			},
			"SubnetArn": &graphql.Field{
				Type: graphql.String,
			},
			"SubnetId": &graphql.Field{
				Type: graphql.String,
			},
			"Tags": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
			"VpcId": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var KarnaGraphQLSecurityGroupType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "SecurityGroup",
		Fields: graphql.Fields{
			"Description": &graphql.Field{
				Type: graphql.String,
			},
			"GroupId": &graphql.Field{
				Type: graphql.String,
			},
			"GroupName": &graphql.Field{
				Type: graphql.String,
			},
			"IpPermissions": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
			"IpPermissionsEgress": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
			"OwnerId": &graphql.Field{
				Type: graphql.String,
			},
			"Tags": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
			"VpcId": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var KarnaGraphQLInstanceType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Instance",
		Fields: graphql.Fields{
			"AmiLaunchIndex": &graphql.Field{
				Type: graphql.Int,
			},
			"BlockDeviceMappings": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
			"ImageId": &graphql.Field{
				Type: graphql.String,
			},
			"InstanceId": &graphql.Field{
				Type: graphql.String,
			},
			"KeyName": &graphql.Field{
				Type: graphql.String,
			},
			"LaunchTime": {
				Type: graphql.String,
			},
			"NetworkInterfaces": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
			"PrivateDnsName": {
				Type: graphql.String,
			},
			"PrivateIpAddress": {
				Type: graphql.String,
			},
			"PublicDnsName": {
				Type: graphql.String,
			},
			"PublicIpAddress": {
				Type: graphql.String,
			},
			"SecurityGroups": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
			"SourceDestCheck": &graphql.Field{
				Type: graphql.Int,
			},
			"SubnetId": {
				Type: graphql.String,
			},
			"Tags": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
			"VpcId": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
