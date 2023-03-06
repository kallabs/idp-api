package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/kallabs/sso-api/src/internal/infra"
	"github.com/kallabs/sso-api/src/internal/interfaces"
	"github.com/kallabs/sso-api/src/internal/interfaces/repos"
	"github.com/kallabs/sso-api/src/internal/interfaces/services"
	"github.com/kallabs/sso-api/src/internal/utils"
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

	db, err := infra.NewMariaDB(utils.Conf.DatabaseUri)
	if err != nil {
		fmt.Print(err)
		return
	}
	repos := repos.NewRepos(db)
	fmt.Printf("User API server listening %s\n", utils.Conf.ServerAddress)

	services := services.NewServices(repos)

	srv, err := interfaces.NewHTTPServer(utils.Conf.ServerAddress, repos, services)

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
