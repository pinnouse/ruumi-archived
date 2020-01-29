package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"time"
)

func newSession() *s3.S3 {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess, &aws.Config{
		Endpoint: aws.String(os.Getenv("S3_ENDPOINT")),
	})
	return svc
}

func getObject(svc *s3.S3, objectKey string) (urlStr string, err error) {
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET")),
		Key:    aws.String(objectKey),
	})
	urlStr, err = req.Presign(10 * time.Hour)
	return
}
