package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gplume/no-mux/middle"
)

var logger *log.Logger

func main() {
	logger = log.New(os.Stdout, "server: ", log.Lshortfile)

	api := &API{
		// route: "/"
		HomeHandler: middle.Ware(new(HomeHandler),
			// with:
			middle.Logging(logger),
			middle.Notify(logger),
		),
		// route: "/user"
		UserHandler: middle.Ware(new(UserHandler),
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

// https://blog.merovius.de/2017/06/18/how-not-to-use-an-http-router.html
// https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81
// https://husobee.github.io/golang/url-router/2015/06/15/why-do-all-golang-url-routers-suck.html

// CutPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
func CutPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

// API ...
type API struct {
	HomeHandler http.Handler
	UserHandler http.Handler
}

func (h *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("******************************")
	fmt.Println("API ServeHTTP Method called")
	fmt.Println(CutPath(r.URL.Path))
	fmt.Println("******************************")
	var head string
	head, r.URL.Path = CutPath(r.URL.Path)
	switch head {
	case "user":
		h.UserHandler.ServeHTTP(w, r)
		return
	case "":
		h.HomeHandler.ServeHTTP(w, r)
		return
	}
	http.Error(w, fmt.Sprintf("Path: %q Not Found", r.URL.Path), http.StatusNotFound)
}

// JSMAP shortcut for map[string]interface{}
type JSMAP map[string]interface{}

// JSON ...
func JSON(w http.ResponseWriter, status int, value interface{}) {
	body, err := json.Marshal(value)
	if err != nil {
		logger.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF8")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(status)
	_, err = w.Write(body)
}
