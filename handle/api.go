package handle

import (
	"fmt"
	"net/http"

	"github.com/gplume/no-mux/utils"
)

// API ...
type API struct {
	HomeHandler http.Handler
	UserHandler http.Handler
}

func (h *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("******************************")
	fmt.Println("API ServeHTTP Method called")
	fmt.Println(utils.CutPath(r.URL.Path))
	fmt.Println("******************************")
	var head string
	head, r.URL.Path = utils.CutPath(r.URL.Path)
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
