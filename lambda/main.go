package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Name string `json:"name"`
}

func HandleRequest(ctx context.Context, event MyEvent) (string, error) {
	// https://gobyexample.com/environment-variables
	fmt.Println("tablename:", os.Getenv("tablename"))
	fmt.Println("lettuce:", os.Getenv("lettuce"))
	return fmt.Sprintf("Hello %s!", event.Name), nil
}

func main() {
	lambda.Start(HandleRequest)
}
