package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, event interface{}) (string, error) {
	eventJSON, _ := json.Marshal(event)
	log.Println("event = ", string(eventJSON))
	return "Hello I'm myFunc!", nil
}

func main() {
	lambda.Start(handleRequest)
}
