package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

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

func handler(ctx context.Context, event S3ObjectLambdaEvent) error {
	// Your code
	fmt.Println(event)
	return nil
}

func cmd() {
	lambda.Start(handler)
}
