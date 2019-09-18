package core

import (
	"regexp"
	"strings"
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
