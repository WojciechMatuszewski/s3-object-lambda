package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	lambda.Start(handler)
}

type S3ObjectLambdaEvent struct {
	Xamzrequestid    string `json:"xAmzRequestId"`
	Getobjectcontext struct {
		Inputs3URL  string `json:"inputS3Url"`
		Outputroute string `json:"outputRoute"`
		Outputtoken string `json:"outputToken"`
	} `json:"getObjectContext"`
	Configuration struct {
		Accesspointarn           string `json:"accessPointArn"`
		Supportingaccesspointarn string `json:"supportingAccessPointArn"`
		Payload                  string `json:"payload"`
	} `json:"configuration"`
	Userrequest struct {
		URL     string `json:"url"`
		Headers struct {
			Host              string `json:"Host"`
			AcceptEncoding    string `json:"Accept-Encoding"`
			XAmzContentSha256 string `json:"X-Amz-Content-SHA256"`
		} `json:"headers"`
	} `json:"userRequest"`
	Useridentity struct {
		Type        string `json:"type"`
		Principalid string `json:"principalId"`
		Arn         string `json:"arn"`
		Accountid   string `json:"accountId"`
		Accesskeyid string `json:"accessKeyId"`
	} `json:"userIdentity"`
	Protocolversion string `json:"protocolVersion"`
}

type Response struct {
	StatusCode int `json:"status_code"`
}

func handler(ctx context.Context, event S3ObjectLambdaEvent) (Response, error) {
	spew.Dump(event)

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	s3Url := event.Getobjectcontext.Inputs3URL
	resp, err := http.Get(s3Url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	svc := s3.NewFromConfig(cfg)

	_, err = svc.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key:    aws.String("cat.jpeg"),
	})
	if err != nil {
		panic(err)
	}

	_, err = svc.WriteGetObjectResponse(ctx, &s3.WriteGetObjectResponseInput{
		RequestRoute: &event.Getobjectcontext.Outputroute,
		RequestToken: &event.Getobjectcontext.Outputtoken,
		Body:         resp.Body,
	})
	if err != nil {
		fmt.Println("Failed to send the write get object response")
		panic(err)
	}

	return Response{StatusCode: 200}, nil
}
