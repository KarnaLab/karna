package core

var Lambda KarnaLambdas
var AGW KarnaAPIGateway
var EC2 KarnaEC2
var S3 KarnaS3

func init() {
	EC2.init()
	AGW.init()
	Lambda.init()
	S3.init()
}
