package awsutils

import (
	"bytes"
	"os"

	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func S3PutObject(body []byte, bucket string, key string) {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION"))},
	)
	if err != nil {
		panic(err)
	}

	s3uploader := s3manager.NewUploader(sess)

	_, err = s3uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key: aws.String(key),
		Body: bytes.NewReader(body),
	})
	if err != nil {
		panic(err)
	}

}