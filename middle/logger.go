package middle

import (
	"log"
	"net/http"
)

// Logging ...
func Logging(logger *log.Logger) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Println(">>>>>>>>>>>> Logging() before")
			defer logger.Println("<<<<<<<<<<<< Logging() after")
			h.ServeHTTP(w, r)
		})
	}
}
