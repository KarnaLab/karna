package core

import (
	"context"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func (karnaS3 *KarnaS3) init() {
	cfg, err := external.LoadDefaultAWSConfig()

	if err != nil {
		LogErrorMessage("unable to load SDK config, " + err.Error())
		os.Exit(2)
	}

	karnaS3.Client = s3.New(cfg)
}

func (karnaS3 *KarnaS3) Upload(deployment *KarnaDeployment, archivePath string) (err error) {
	part, _ := ioutil.ReadFile(archivePath)
	input := &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(strings.NewReader(string(part))),
		Key:    aws.String(deployment.File),
		Bucket: aws.String(deployment.Bucket),
	}

	req := karnaS3.Client.PutObjectRequest(input)

	_, err = req.Send(context.Background())

	return
}
