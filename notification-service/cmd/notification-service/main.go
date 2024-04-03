package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"notification-service/internal/config"
	"notification-service/internal/controllers"
	"notification-service/internal/server"
	"notification-service/internal/service"
	"os"
	"os/signal"
	"sync"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, config.Config); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, getenv func(string) string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Kill)
	defer cancel()

	db, err := newDBClient()
	if err != nil {
		return err
	}

	receiverService := service.NewReceiverService(db, getenv("TABLE_NAME"))
	receiverController := controllers.NewReceiverController(receiverService)

	sns, err := newSNSClient()
	if err != nil {
		return err
	}

	notificationService := service.NewNotificationService(sns, getenv("SNS_TOPIC"))
	notificationController := controllers.NewNotificationController(receiverService, notificationService)

	srv := server.NewServer(receiverController, notificationController)

	httpServer := &http.Server{
		Addr:    net.JoinHostPort("0.0.0.0", getenv("API_PORT")),
		Handler: srv,
	}

	go func() {
		log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()

	return nil
}

func newDBClient() (*dynamodb.Client, error) {
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	return dynamodb.NewFromConfig(cfg), nil
}

func newSNSClient() (*sns.Client, error) {
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	return sns.NewFromConfig(cfg), nil
}
