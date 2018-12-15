package handle

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gplume/no-mux/utils"
)

// Home ...
type Home struct{}

// HomeHandler implements ServeHTPP to return an http.Handler (interface satisfying)
func (h *Home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=======================================")
	fmt.Println("|---------WELCOME TO THE API!---------|")
	fmt.Println("=======================================")

	time.Sleep(44 * time.Millisecond)

	utils.JSON(w, http.StatusOK, utils.JSMAP{
		"msg": "OK",
	})
}
