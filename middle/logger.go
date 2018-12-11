package middle

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Logging ...
func Logging(logger *log.Logger) Wrapper {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func(begin time.Time) {
				logger.Println(
					"method", fmt.Sprintf("%T", h),
					"input", r.URL.Query(),
					"took", time.Since(begin).Round(time.Millisecond),
				)
			}(time.Now())
			h.ServeHTTP(w, r)
		})
	}
}
