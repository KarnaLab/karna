package core

var Lambda KarnaLambdas
var AGW KarnaAPIGateway
var EC2 KarnaEC2

func init() {
	EC2.init()
	EC2.BuildEC2Tree()
}
