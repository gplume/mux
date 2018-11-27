package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

// UserHandler ...
type UserHandler struct {
	ProfileHandler *ProfileHandler
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("====================================")
	fmt.Println("UserHandler.ServeHTTP Method called")
	fmt.Println(CutPath(r.URL.Path))
	fmt.Println("====================================")

	// using context to pass some value found in query:
	headUser := r.URL.Query().Get("id")
	ctx := context.WithValue(r.Context(), contextKey("XX--ID"), headUser)

	var head string
	head, r.URL.Path = CutPath(r.URL.Path)
	idArr := strings.Split(head, ",")
	ids := make([]int, len(idArr))
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
		head, _ := CutPath(r.URL.Path)
		switch head {
		case "profile":
			// We can't just make ProfileHandler an http.Handler; it needs the
			// user id. Let's instead… !
			h.ProfileHandler.Handler(headUser).ServeHTTP(w, r)
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

func (h *UserHandler) handleGet(w http.ResponseWriter, r *http.Request, ids ...int) {
	fmt.Println("XX--ID:", r.Context().Value(contextKey("XX--ID")))
	q1 := r.URL.Query().Get("id")

	jsm := JSMAP{}
	for ki, vi := range ids {
		jsm[strconv.Itoa(ki)] = vi
	}
	json := JSMAP{
		"array": jsm,
	}
	if q1 != "" {
		json["query_id"] = q1
	}
	JSON(w, http.StatusOK, json)
}

func (h *UserHandler) handlePut(w http.ResponseWriter, r *http.Request, id ...int) {
	fmt.Printf("handlePut ID: %v", id)
	fmt.Println("")
}

// ProfileHandler .
/////////////////////
type ProfileHandler struct {
}

// Handler .
func (h *ProfileHandler) Handler(id string) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Do whatever
	})
}
