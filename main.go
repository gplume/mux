package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"time"
)

// https://blog.merovius.de/2017/06/18/how-not-to-use-an-http-router.html
// https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81

// ShiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

var logger *log.Logger

func main() {
	var wait time.Duration
	logger = log.New(os.Stdout, "server: ", log.Lshortfile)

	api := &API{
		HomeHandler: Adapt(new(HomeHandler),
			// with middlewares:
			// RecoverFromPanic(logger),
			Notify(logger),
			Logging(logger),
		),
		UserHandler: Adapt(new(UserHandler),
			// with middlewares:
			// RecoverFromPanic(logger),
			Notify(logger),
			Logging(logger),
		),
	}

	srv := &http.Server{
		Addr: ":8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      RecoverFromPanic(logger, api),
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
	ctx, cancel := context.WithTimeout(context.Background(), wait)
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

// API ...
type API struct {
	HomeHandler http.Handler
	UserHandler http.Handler
}

func (h *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("******************************")
	fmt.Println("API ServeHTTP Method called")
	fmt.Println(ShiftPath(r.URL.Path))
	fmt.Println("-------------------------")
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	switch head {
	case "user":
		h.UserHandler.ServeHTTP(w, r)
		return
	case "":
		h.HomeHandler.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Not Found", http.StatusNotFound)
}
