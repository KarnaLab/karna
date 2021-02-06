package deploy

import (
	"context"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

type KarnaS3Model struct {
	Client *s3.Client
}

func (KarnaS3Model *KarnaS3Model) init() {
	cfg, err := external.LoadDefaultAWSConfig()

	if err != nil {
		logger.Error("unable to load SDK config, " + err.Error())
	}

	KarnaS3Model.Client = s3.New(cfg)
}

//Upload => Will upload to S3 specified archive.
func (KarnaS3Model *KarnaS3Model) upload(deployment *KarnaDeployment, archivePath string) (err error) {
	part, _ := ioutil.ReadFile(archivePath)
	input := &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(strings.NewReader(string(part))),
		Key:    aws.String(deployment.File),
		Bucket: aws.String(deployment.Bucket),
	}

	req := KarnaS3Model.Client.PutObjectRequest(input)

	_, err = req.Send(context.Background())

	return
}
