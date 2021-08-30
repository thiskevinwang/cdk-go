package main

import (
	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/awss3assets"
	"github.com/aws/aws-cdk-go/awscdk/awssns"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type CdkLambdaGoStackProps struct {
	awscdk.StackProps
}

func NewCdkLambdaGoStack(scope constructs.Construct, id string, props *CdkLambdaGoStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// The code that defines your stack goes here

	// as an example, here's how you would define an AWS SNS topic:
	awssns.NewTopic(stack, jsii.String("MyTopic"), &awssns.TopicProps{
		DisplayName: jsii.String("MyCoolTopic"),
	})

	// Adding a basic Dynamo DB table
	table := awsdynamodb.NewTable(stack, jsii.String("MyTable"), &awsdynamodb.TableProps{
		TableName: jsii.String("GoCDKTable"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("pk"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: jsii.String("sk"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		BillingMode: awsdynamodb.BillingMode_PAY_PER_REQUEST,
	})

	// RUN
	// `GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./lambda/main ./lambda/main.go`
	// https://stackoverflow.com/questions/58133166/getting-error-fork-exec-var-task-main-no-such-file-or-directory-while-execut
	awslambda.NewFunction(stack, jsii.String("MyGoFunction"), &awslambda.FunctionProps{
		Description: jsii.String("A lambda function, written in Go"),
		Runtime:     awslambda.Runtime_GO_1_X(),
		// In the `lambda` folder, there needs to be an *executable*, preferably named `main`
		Code: awslambda.Code_FromAsset(jsii.String("lambda"), &awss3assets.AssetOptions{}),
		// In the original `go` code from `lambda/main.go`, there needs to be a func defined:
		// - `func main()`...
		Handler: jsii.String("main"),
	})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	stackId := "GolangStack"

	NewCdkLambdaGoStack(app, stackId, &CdkLambdaGoStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
