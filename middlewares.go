package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// RecoverFromPanic as main recover for the global Handler...
func RecoverFromPanic(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Println("err RecoverFromPanic2", rec,
					"http.url", r.RequestURI, "http.path", r.URL.Path,
					"http.method", r.Method, "http.user_agent", r.Header.Get("User-Agent"),
					"http.proto", r.Proto)
				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error": fmt.Sprintf("%v, %T", rec, rec),
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Adapt Our Adapt function takes the handler you want to adapt,
// and a list of our Adapter types. The result of our Notify function
// is an acceptable Adapter. Our Adapt function will simply iterate over all adapters,
//  calling them one by one (in reverse order) in a chained manner, returning the result of the first adapter.
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	// reverse order:
	// for _, adapter := range adapters {
	// 	h = adapter(h)
	// }

	// straight order:
	for i := len(adapters) - 1; i >= 0; i-- {
		h = adapters[i](h)
	}
	return h
}

// Adapter type (it gets its name from the adapter pattern — also known as the decorator pattern)
// above is a function that both takes in and returns an http.Handler. This is the essence of the wrapper;
//we will pass in an existing http.Handler, the Adapter will adapt it, and return a new (probably wrapped) http.Handler
// for us to use in its place. So far this is not much different from just wrapping http.HandlerFunc types,
// however, now, we can instead write functions that themselves return an Adapter.
type Adapter func(http.Handler) http.Handler

// Notify ...
func Notify(logger *log.Logger) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Println("Notify() before")
			defer func(begin time.Time) {
				logger.Printf("Notify() after %vs", time.Since(begin).Seconds())
			}(time.Now())
			h.ServeHTTP(w, r)
		})
	}
}

// Logging ...
func Logging(logger *log.Logger) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Println("Logging() before")
			defer logger.Println("Logging() after")
			h.ServeHTTP(w, r)
		})
	}
}
