package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v2.0/computervision"
	"github.com/Azure/go-autorest/autorest"
)

var computerVisionContext context.Context

func main() {
	computerVisionKey := os.Getenv("COMPUTER_VISION_SUBSCRIPTION_KEY")
	endPointURL := os.Getenv("COMPUTER_VISION_ENDPOINT")

	computerVisionClient := computervision.New(endPointURL)
	computerVisionClient.Authorizer = autorest.NewCognitiveServicesAuthorizer(computerVisionKey)

	computerVisionContext = context.Background()

	url := "/Users/shinnosuke/Desktop/ocr/IMG_2561.JPG"
	image := computervision.ImageURL{URL: &url}

	textHeaders, err := computerVisionClient.BatchReadFile(computerVisionContext, image)
	if err != nil {
		log.Fatal(err)
	}

	operationLocation := autorest.ExtractHeaderValue("Operation-Location", textHeaders.Response)

	numberOfCharsInOperationId := 36
	operationId := string(operationLocation[len(operationLocation)-numberOfCharsInOperationId : len(operationLocation)])

	readOperationResult, err := computerVisionClient.GetReadOperationResult(computerVisionContext, operationId)
	if err != nil {
		log.Fatal(err)
	}

	i := 0
	maxRetries := 10

	fmt.Println("Recognizing text in a remote image with the batch Read API ...")
	for readOperationResult.Status != computervision.Failed && readOperationResult.Status != computervision.Succeeded {
		if i >= maxRetries {
			break
		}
		i++

		fmt.Printf("Server status: %v, waiting %v seconds...\n", readOperationResult.Status, i)
		time.Sleep(1 * time.Second)

		readOperationResult, err = computerVisionClient.GetReadOperationResult(computerVisionContext, operationId)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println()
	for _, recResult := range *(readOperationResult.RecognitionResults) {
		for _, line := range *recResult.Lines {
			fmt.Println(*line.Text)
		}
	}
}
