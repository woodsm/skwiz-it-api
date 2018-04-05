package storage

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/developmentnow/skwiz-it-api/config"
	"github.com/twinj/uuid"
	"log"
)

var conf = config.LoadConfig()

func SaveToS3(b64 string) (string, error) {
	cred := credentials.NewStaticCredentials(conf.S3.AccessKey, conf.S3.AccessSecret, "")
	_, err := cred.Get()
	if err != nil {
		fmt.Printf("bad credentials: %s", err)
		return "", err
	}

	cfg := aws.NewConfig().WithRegion(conf.S3.Region).WithCredentials(cred)
	svc := s3.New(session.New(), cfg)

	byteArray, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		fmt.Printf("Unable to read base64 into byte array: %s", err)
		return "", err
	}

	fileId := uuid.NewV4().String() + ".png"

	fileBytes := bytes.NewReader(byteArray)
	fileType := http.DetectContentType(byteArray)
	acl := "public-read"
	params := &s3.PutObjectInput{
		Bucket: aws.String(conf.S3.Bucket),
		Key:    aws.String(fileId),
		Body:   fileBytes,
		// ContentLength: aws.Int64(size),
		ContentType: aws.String(fileType),
		ACL:         &acl,
	}

	resp, err := svc.PutObject(params)
	if err != nil {
		fmt.Printf("bad response: %s", err)
		return "", err
	}

	log.Printf("Uploaded base64 image %s to S3 identified as %s", fileId, resp)
	return fileId, nil
}
