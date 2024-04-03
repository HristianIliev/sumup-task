package main

import (
	"context"
	"log"
	"strconv"
	"sync"

	"notification-processor/internal/config"
	"notification-processor/internal/sqsclient"
	"notification-processor/internal/sqsclient/subscriber"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal("error loading AWS config:", err)
		return
	}

	client := sqs.NewFromConfig(cfg)

	subscribers := []sqsclient.NotificationSubscriber{
		&subscriber.EmailSubscriber{},
		&subscriber.SMSSubscriber{},
		&subscriber.SlackSubscriber{},
	}

	var wg sync.WaitGroup

	workers, err := strconv.Atoi(config.Config("WORKERS"))
	if err != nil {
		log.Fatalf("number of workers configuration is not a valid int")
		return
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)

		sqsClient := sqsclient.New(client, config.Config("QUEUE_URL"))

		pollDelay, err := strconv.Atoi(config.Config("POLL_DELAY"))
		if err != nil {
			log.Fatalf("poll delay configuration is not a valid int")
			return
		}

		consumer := sqsclient.NewConsumer(i, sqsClient, pollDelay, subscribers)

		go consumer.Run(ctx)
	}

	wg.Wait()
}
