package core

// Lambda : Lambda Model for Karna core.
var Lambda KarnaLambdaModel

// AGW : APIGateway Model for Karna core.
var AGW KarnaAPIGatewayModel

// EC2 : EC2 Model for Karna core.
var EC2 KarnaEC2Model

// S3 : S3 Model for Karna core.
var S3 KarnaS3Model
var logger *KarnaLogger

func init() {
	EC2.init()
	AGW.init()
	Lambda.init()
	S3.init()
}
