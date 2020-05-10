package awslib

import (
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//AwsManager holds info for the aws manager
type AwsManager struct {
	BucketName       string
	Region           string
	ExtensionMapping map[string]string
	Session          *session.Session
}

var awsManager *AwsManager

//GetRegion returns the geo-location of the host
func GetRegion() string {
	return awsManager.Region
}

//GetBucket returns the bucketName
func GetBucket() string {
	return awsManager.BucketName
}

//GetFileExtension returns file extension of the queried fileName
func GetFileExtension(filename string) string {
	index := strings.LastIndex(filename, ".")
	fileExtension := "binary/octet-stream"
	if index >= 0 { //valid file names: a.png, .jpg, etc
		ext := "." + filename[index+1:]
		fileExtension = awsManager.ExtensionMapping[ext]
	}
	return fileExtension
}

//CreateBucketInS3 creates a new bucket in s3
func CreateBucketInS3(bucketName string) error {
	svc := s3.New(awsManager.Session)

	_, err := svc.CreateBucket(&s3.CreateBucketInput{Bucket: &bucketName})
	if err != nil {
		return fmt.Errorf("failed to create new bucket, err: %v", err)
	}

	return nil
}

//UploadToS3 uploads data to s3
func UploadToS3(key string, body io.Reader, metaData map[string]*string, bucketName string) error {
	uploader := s3manager.NewUploader(awsManager.Session)
	//ACL logic to be accumulated here
	if bucketName == "" {
		bucketName = awsManager.BucketName
	}
	_, err := uploader.Upload(&s3manager.UploadInput{
		Body:        body,
		Bucket:      aws.String(bucketName),
		ACL:         aws.String("public-read"),
		Key:         aws.String(key),
		Metadata:    metaData,
		ContentType: aws.String(GetFileExtension(key)),
	})
	if err != nil {
		return err
	}
	return nil
}
