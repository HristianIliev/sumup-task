package main

import (
	"context"
	"ethereum-fetcher/internal/config"
	"ethereum-fetcher/internal/controllers"
	"ethereum-fetcher/internal/server"
	"ethereum-fetcher/internal/service/transaction"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
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

	db, err := initializeDB(getenv("DB_CONNECTION_URL"))
	if err != nil {
		return err
	}

	transactionService := transaction.New(db)
	transactionController := controllers.New(*transactionService)
	srv := server.NewServer(transactionController)

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

func initializeDB(connectionURL string) (*sqlx.DB, error) {
	log.Println(connectionURL)
	db, err := sqlx.Connect("pgx", connectionURL)
	if err != nil {
		log.Fatalln(err)
	}

	return db, nil
}
