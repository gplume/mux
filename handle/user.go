package handle

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gplume/no-mux/utils"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

// User ...
type User struct {
	Profile *ProfileHandler
}

func (h *User) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("====================================")
	fmt.Println("UserHandler.ServeHTTP Method called")
	fmt.Println(utils.CutPath(r.URL.Path))
	fmt.Println("====================================")

	// using context to pass some value found in query:
	userID := r.URL.Query().Get("id")
	ctx := context.WithValue(r.Context(), contextKey("XX--ID"), userID)

	var head string
	head, r.URL.Path = utils.CutPath(r.URL.Path)
	idArr := strings.Split(head, ",")
	ids := make([]int, 0, len(idArr))
	if len(idArr) >= 1 {
		for _, v := range idArr {
			id, err := strconv.Atoi(v)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid user id %q", head), http.StatusBadRequest)
				return
			}
			ids = append(ids, id)
		}
	} else {
		http.Error(w, fmt.Sprintf("Invalid user id %q", head), http.StatusBadRequest)
	}
	if r.URL.Path == "/" {
		switch r.Method {
		case "GET":
			h.handleGet(w, r.WithContext(ctx), ids...)
		case "PUT":
			h.handlePut(w, r, ids...)
		default:
			http.Error(w, "Only GET and PUT are allowed", http.StatusMethodNotAllowed)
		}
	} else {
		head, _ := utils.CutPath(r.URL.Path)
		switch head {
		case "profile":
			// We can't just make ProfileHandler an http.Handler; it needs the
			// user id. Let's instead… !
			h.Profile.Handler(userID).ServeHTTP(w, r)
		case "account":
			// Left as an exercise to the reader.
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
		return
	}
	if r.URL.Path == "/extends" {
		http.Error(w, "Not Implemented yet", http.StatusForbidden)
	}
}

func (h *User) handleGet(w http.ResponseWriter, r *http.Request, ids ...int) {
	fmt.Println("XX--ID:", r.Context().Value(contextKey("XX--ID")))
	q1 := r.URL.Query().Get("id")

	jsm := utils.JSMAP{}
	for ki, vi := range ids {
		jsm[strconv.Itoa(ki)] = vi
	}
	json := utils.JSMAP{
		"array": jsm,
	}
	if q1 != "" {
		json["query_id"] = q1
	}
	utils.JSON(w, http.StatusOK, json)
}

func (h *User) handlePut(w http.ResponseWriter, r *http.Request, id ...int) {
	fmt.Printf("handlePut ID: %v", id)
	fmt.Println("")
}

// ProfileHandler .
/////////////////////
type ProfileHandler struct {
}

// Handler .
func (h *ProfileHandler) Handler(id string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, rs *http.Request) {
		// Do whatever
		utils.JSON(w, 200, utils.JSMAP{
			"PROFILEhandler": "OK",
		})
	})
}
