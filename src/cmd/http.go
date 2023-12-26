package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/kallabs/idp-api/src/internal/adapters"
	"github.com/kallabs/idp-api/src/internal/adapters/storage"
	"github.com/kallabs/idp-api/src/internal/utils"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	err := utils.LoadConfig()
	if err != nil {
		fmt.Print(err)
		return
	}

	db, err := storage.ConnectPostgres(utils.Conf.DatabaseUri)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Printf("User API server listening %s\n", utils.Conf.ServerAddress)

	srv, err := adapters.NewHTTPServer(utils.Conf.ServerAddress, db)

	if err != nil {
		log.Fatal(err)
		return
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
