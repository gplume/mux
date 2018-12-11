package handle

import (
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/gplume/no-mux/logger"
	"github.com/gplume/no-mux/utils"
)

// Home ...
type Home struct{}

// HomeHandler implements ServeHTPP to return an http.Handler (interface satisfying)
func (h *Home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=======================================")
	fmt.Println("|---------WELCOME TO THE API!---------|")
	fmt.Println("=======================================")

	q1 := r.URL.Query().Get("id")
	logger.Log.Println("ID:", q1, reflect.TypeOf(q1))

	q2 := r.URL.Query().Get("page")
	logger.Log.Println("PAGE:", q2, reflect.TypeOf(q2))

	time.Sleep(444 * time.Millisecond)

	utils.JSON(w, http.StatusOK, utils.JSMAP{
		"msg": "OK",
	})
}
