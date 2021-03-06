package main

import (
	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/awseks"
	"github.com/aws/aws-cdk-go/awscdk/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/awslambdago"
	"github.com/aws/aws-cdk-go/awscdk/awssns"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type StackProps struct {
	awscdk.StackProps
}

const baseName = "GoCDK"

func Stack(scope constructs.Construct, id string, props *StackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)
	stack.StackId()

	// The code that defines your stack goes here

	// as an example, here's how you would define an AWS SNS topic:
	awssns.NewTopic(stack, jsii.String(baseName+"SnsTopic"), &awssns.TopicProps{})

	// Adding a basic Dynamo DB table
	table := awsdynamodb.NewTable(stack, jsii.String(baseName+"DynamoTable"), &awsdynamodb.TableProps{
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

	fn := awslambdago.NewGoFunction(stack, jsii.String(baseName+"Lambda"), &awslambdago.GoFunctionProps{
		Description: jsii.String("A lambda function, written in Go"),
		Runtime:     awslambda.Runtime_GO_1_X(),
		Entry:       jsii.String("lambda/main.go"),
		Environment: &map[string]*string{"tablename": jsii.String(*table.TableName()), "lettuce": jsii.String("69")},
	})

	table.GrantReadWriteData(fn)

	awseks.NewFargateCluster(stack, jsii.String(baseName+"FargateCluster"), &awseks.FargateClusterProps{
		Version: awseks.KubernetesVersion_V1_21(),
	})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	stackId := "GolangStack"

	Stack(app, stackId, &StackProps{
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
