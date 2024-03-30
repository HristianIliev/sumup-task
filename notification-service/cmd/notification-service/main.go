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

	dynamodbconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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

	db, err := newclient()
	if err != nil {
		return err
	}

	receiverService := service.NewReceiverService(db, getenv("TABLE_NAME"))
	receiverController := controllers.New(receiverService)
	srv := server.NewServer(receiverController)

	httpServer := &http.Server{
		Addr:    net.JoinHostPort("0.0.0.0", getenv("API_PORT")),
		Handler: srv,
	}

	go func() {
		log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
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
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()

	return nil
}

func newclient() (*dynamodb.Client, error) {
	cfg, err := dynamodbconfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	return dynamodb.NewFromConfig(cfg), nil
}
