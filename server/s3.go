package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Client ...
type S3Client struct {
	Client *s3.S3
}

// NewS3Client ...
func NewS3Client(accessKey, secretKey, region string) *S3Client {
	token := ""
	creds := credentials.NewStaticCredentials(accessKey, secretKey, token)
	_, err := creds.Get()
	if err != nil {
		fmt.Printf("bad credentials: %s", err)
	}
	cfg := aws.NewConfig().WithRegion("us-east-1").WithCredentials(creds)
	svc := s3.New(session.New(), cfg)
	return &S3Client{
		Client: svc,
	}
}
