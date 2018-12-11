package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gplume/no-mux/handle"
	"github.com/gplume/no-mux/logger"
	"github.com/gplume/no-mux/middle"
)

func main() {

	logger := logger.New()
	api := &handle.API{
		// route: "/"
		HomeHandler: middle.Ware(new(handle.Home),
			// with:
			middle.Logging(logger),
			middle.Notify(logger),
		),
		// route: "/user"
		UserHandler: middle.Ware(new(handle.User),
			// with:
			middle.Logging(logger),
			middle.Notify(logger),
		),
	}

	srv := &http.Server{
		Addr: ":8080",
		// It is good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      middle.RecoverFromPanic(logger, api),
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)
	logger.Println("---ready---")
	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	logger.Println("shutting down")
	os.Exit(0)
}
