package core

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (karnaS3 *KarnaS3) init() {
	cfg, err := external.LoadDefaultAWSConfig()

	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	karnaS3.Client = s3.New(cfg)
}

func (karnaS3 *KarnaS3) Upload(path) {
	input := &s3.UploadPartInput{}

	req := karnaS3.Client.UploadPartRequest(input)

	response, err := req.Send(context.Background())

	if err != nil {
		// Abort upload.
		panic(err.Error())
	}

	fmt.Println(response)
}
