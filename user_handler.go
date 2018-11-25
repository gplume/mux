package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

// UserHandler ...
type UserHandler struct {
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UserHandler.ServeHTTP Method called")
	fmt.Println(ShiftPath(r.URL.Path))
	fmt.Println("-------------------------")
	var head, tail string
	head, tail = ShiftPath(r.URL.Path)
	id, err := strconv.Atoi(head)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid user id %q", head), http.StatusBadRequest)
		return
	}
	if tail == "/" {
		switch r.Method {
		case "GET":
			h.handleGet(id, r)
		case "PUT":
			h.handlePut(id)
		default:
			http.Error(w, "Only GET and PUT are allowed", http.StatusMethodNotAllowed)
		}
	}
	if tail == "/extends" {
		http.Error(w, "Not Implemented yet", http.StatusForbidden)
	}
}

func (h *UserHandler) handleGet(id int, r *http.Request) {
	fmt.Printf("handleGet ID: %v", id)
	fmt.Println("")

	q1 := r.URL.Query().Get("id")
	logger.Println("Query ID:", q1, reflect.TypeOf(q1))
}

func (h *UserHandler) handlePut(id int) {
	fmt.Printf("handlePut ID: %v", id)
	fmt.Println("")
}
