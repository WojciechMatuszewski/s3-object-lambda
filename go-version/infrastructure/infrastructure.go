package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/awss3"
	"github.com/aws/aws-cdk-go/awscdk/awss3assets"
	"github.com/aws/aws-cdk-go/awscdk/awss3objectlambda"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type AppProps struct {
	awscdk.StackProps
}

func NewApp(scope constructs.Construct, id string, props *AppProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	filesBucket := awss3.NewBucket(stack, jsii.String("filesBucket"), &awss3.BucketProps{
		AccessControl:     awss3.BucketAccessControl_BUCKET_OWNER_FULL_CONTROL,
		BlockPublicAccess: awss3.BlockPublicAccess_BLOCK_ALL(),
		AutoDeleteObjects: jsii.Bool(true),
		RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
	})

	filesBucketAccessPoint := awss3.NewCfnAccessPoint(stack, jsii.String("filesBucketAccessPoint"), &awss3.CfnAccessPointProps{
		Bucket: filesBucket.BucketName(),
		Name:   jsii.String("accesspoint"),
	})

	oneTimeReceiverLambda := awslambda.NewFunction(stack, jsii.String("objectTransformer"), &awslambda.FunctionProps{
		Code: awslambda.AssetCode_FromAsset(jsii.String("src"), &awss3assets.AssetOptions{
			AssetHash: jsii.String(lambdaHash()),
		}),
		Handler: jsii.String("main"),
		Runtime: awslambda.Runtime_GO_1_X(),
		Environment: &map[string]*string{
			"BUCKET_NAME": filesBucket.BucketName(),
		},
	})

	oneTimeReceiverLambda.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions: &[]*string{
			jsii.String("s3-object-lambda:WriteGetObjectResponse"),
		},
		Effect: awsiam.Effect_ALLOW,
		Resources: &[]*string{
			jsii.String("*"),
		},
	}))

	oneTimeReceiverLambda.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions: &[]*string{
			jsii.String("s3:DeleteObject"),
		},
		Effect: awsiam.Effect_ALLOW,
		Resources: &[]*string{
			filesBucket.ArnForObjects(jsii.String("*")),
		},
	}))

	awss3objectlambda.NewCfnAccessPoint(stack, jsii.String("objectLambdaAccessPoint"), &awss3objectlambda.CfnAccessPointProps{
		Name: jsii.String("oblambdaap"),
		ObjectLambdaConfiguration: awss3objectlambda.CfnAccessPoint_ObjectLambdaConfigurationProperty{
			SupportingAccessPoint: awscdk.Fn_Sub(
				jsii.String("arn:${PARTITION}:s3:${REGION}:${ACCOUNT_ID}:accesspoint/${BUCKET_ACCESS_POINT}"),
				&map[string]*string{
					"PARTITION":           awscdk.Aws_PARTITION(),
					"REGION":              awscdk.Aws_REGION(),
					"ACCOUNT_ID":          awscdk.Aws_ACCOUNT_ID(),
					"BUCKET_ACCESS_POINT": awscdk.Fn_Ref(filesBucketAccessPoint.LogicalId()),
				}),
			TransformationConfigurations: []awss3objectlambda.CfnAccessPoint_TransformationConfigurationProperty{
				{
					Actions: &[]*string{
						jsii.String("GetObject"),
					},
					ContentTransformation: map[string]map[string]string{
						"AwsLambda": {
							"FunctionArn": *oneTimeReceiverLambda.FunctionArn(),
						},
					},
				},
			},
			CloudWatchMetricsEnabled: jsii.Bool(true),
		},
	}).AddDependsOn(filesBucketAccessPoint)

	awscdk.NewCfnOutput(stack, jsii.String("filesBucketName"), &awscdk.CfnOutputProps{
		Value: filesBucket.BucketName(),
	})

	awscdk.NewCfnOutput(stack, jsii.String("filesBucketArn"), &awscdk.CfnOutputProps{
		Value: filesBucket.BucketArn(),
	})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewApp(app, "App", &AppProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}

func lambdaHash() string {
	h := sha256.New()
	f, err := os.Open("./src/main")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return hex.EncodeToString(h.Sum(nil))
}
