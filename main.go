package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("Akashgreninja"))
	//https://github.com/aws/aws-sdk-go-v2/blob/main/service/s3/api_op_GetBucketNotificationConfiguration.go

	if err != nil {
		log.Fatal("coming from here")
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)

	//adding events
	//https://github.com/aws/aws-sdk-go-v2/blob/main/service/s3/api_op_PutBucketNotificationConfiguration.go
	//the below code wont wonk for V2 of the golang api

	// notificationConfig := &s3.PutBucketNotificationConfigurationInput{
	// 	Bucket: aws.String("greninja"),
	// 	NotificationConfiguration: &s3.NotificationConfiguration{
	// 		LambdaFunctionConfigurations: []s3.LambdaFunctionConfiguration{
	// 			{
	// 				Events: []s3.Event{
	// 					s3.EventS3ObjectCreated,
	// 					s3.EventS3ObjectRemoved,
	// 				},
	// 				LambdaFunctionArn: aws.String("arn:aws:lambda:ap-south-1:678597231768:function:fordataget"),
	// 			},
	// 		},
	// 	},
	// }

	// _, err = client.PutBucketNotificationConfiguration(notificationConfig)
	// if err != nil {
	// 	log.Fatalf("unable to set up event notification, %v", err)
	// }

	fmt.Println("S3 event notification configured successfully.")

	type result struct {
		Output *s3.PutObjectOutput
		Err    error
	}
	results := make(chan result, 2)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		content, err := ioutil.ReadFile("./8de3439d-60fd-493e-91da-c7aeef70a2d8.jpeg")
		if err != nil {
			results <- result{Err: err}
			return
		}
		output, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String("greninja"),
			Key:    aws.String("pfdfdflus"),
			Body:   strings.NewReader(string(content)),
		})
		results <- result{Output: output, Err: err}
	}()
	go func() {
		defer wg.Done()
		content, err := ioutil.ReadFile("./giving a talk.jpg")
		if err != nil {
			results <- result{Err: err}
			return
		}
		output, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String("greninja"),
			Key:    aws.String("minfdfdus"),
			Body:   strings.NewReader(string(content)),
		})
		results <- result{Output: output, Err: err}
	}()
	wg.Wait()

	close(results)
	for result := range results {
		if result.Err != nil {
			log.Printf("error: %v", result.Err)
			continue
		}
		fmt.Printf("etag: %v", aws.ToString(result.Output.ETag))
	}
}

//https://aws.github.io/aws-sdk-go-v2/docs/getting-started/
