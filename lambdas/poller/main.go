package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context) (string, error) {
	return "hello world", nil
}

func main() {
	lambda.Start(handleRequest)
}
