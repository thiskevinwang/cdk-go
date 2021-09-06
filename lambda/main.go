package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	// https://github.com/aws/aws-sdk-go-v2/issues/1018#issuecomment-755433782
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/jsii-runtime-go"
)

type MyEvent struct {
	Pk string `json:"pk"`
	Sk string `json:"sk"`
}

type MyItem struct {
	Pk   string `json:"pk"`
	Sk   string `json:"sk"`
	Test string `json:"test,omitempty"`
}

func HandleRequest(ctx context.Context, event MyEvent) (string, error) {
	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Using the Config value, create the DynamoDB client
	svc := dynamodb.NewFromConfig(cfg)

	tableName := jsii.String(os.Getenv("tablename"))

	getItemOutput, getItemError := svc.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: tableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: event.Pk},
			"sk": &types.AttributeValueMemberS{Value: event.Sk},
		},
	})
	if getItemError != nil {
		return fmt.Sprintf("getItemError: %s!", getItemError), nil
	}

	// "unmarshallErr unmarshal failed, cannot unmarshal to non-pointer value, got main.MyItem!"
	// var item MyItem

	// "unmarshallErr unmarshal failed, cannot unmarshal to nil value, *main.MyItem!"
	// var item *MyItem

	item := &MyItem{}

	unmarshallErr := attributevalue.UnmarshalMap(getItemOutput.Item, item)
	if unmarshallErr != nil {
		return fmt.Sprintf("unmarshallErr %s!", unmarshallErr), nil
	}
	return fmt.Sprintf("Hello %s!", item.Pk), nil
	// "Hello &{value1 {}}!"

}

func main() {
	lambda.Start(HandleRequest)
}
