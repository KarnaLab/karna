package core

import (
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

func findAPIGatewayID(statement awsPolicyStatement) (apiID string) {
	IDRegex := regexp.MustCompile(`:\w+\/`)
	apiID = IDRegex.FindString(statement.Condition.ArnLike["AWS:SourceArn"])
	apiID = strings.NewReplacer(":", "", "/", "").Replace(apiID)
	return
}

func findS3Bucket(statement awsPolicyStatement) (bucket string) {
	IDRegex := regexp.MustCompile(`::.+`)
	bucket = IDRegex.FindString(statement.Condition.ArnLike["AWS:SourceArn"])
	bucket = strings.NewReplacer(":", "").Replace(bucket)
	return
}

func findCloudWatch(statement awsPolicyStatement) (rule string) {
	IDRegex := regexp.MustCompile(`:rule.+`)
	rule = IDRegex.FindString(statement.Condition.ArnLike["AWS:SourceArn"])
	rule = strings.NewReplacer(":", "").Replace(rule)
	return
}

func findAPIGatewayIDInDomainName(domainName string) (apiID string) {
	IDRegex := regexp.MustCompile(`^:\w+.`)
	apiID = IDRegex.FindString(domainName)
	apiID = strings.NewReplacer(".", "").Replace(apiID)
	return
}

func findService(service string, statements []awsPolicyStatement) (apiGatewayStatement awsPolicyStatement, found bool) {
	for _, statement := range statements {
		if statement.Principal.Service == service {
			found = true
			apiGatewayStatement = statement
			return
		}
		found = false
	}
	return
}

func uniq(slice []string) (uniqValues []string) {
	var uniqMap = make(map[string]bool)

	for _, value := range slice {
		uniqMap[value] = true
	}

	for index, _ := range uniqMap {
		uniqValues = append(uniqValues, index)
	}

	return
}

func findStage(stages []KarnaAGWStage, stage string) (index int) {
	for i, s := range stages {
		if s.Stage == stage {
			index = i
		}
	}
	return
}

func findAlias(aliases []lambda.AliasConfiguration, aliasName string) (alias *lambda.AliasConfiguration) {
	for _, a := range aliases {
		if *a.Name == aliasName {
			alias = &a
		}
	}
	return
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func findInt(needle int, haystack []int) (found bool) {
	for _, value := range haystack {
		if needle == value {
			found = true
		}
	}
	return
}
