package core

var Lambda KarnaLambdas
var AGW KarnaAPIGateway

func init() {
	AGW.init()
	AGW.BuildAGWTree()
}
