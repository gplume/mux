package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// UserHandler ...
type UserHandler struct{}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("====================================")
	fmt.Println("UserHandler.ServeHTTP Method called")
	fmt.Println(CutPath(r.URL.Path))
	fmt.Println("====================================")

	var head string
	head, r.URL.Path = CutPath(r.URL.Path)
	ids := make([]int, 0)
	idArr := strings.Split(head, ",")
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
			h.handleGet(w, r, ids...)
		case "PUT":
			h.handlePut(w, r, ids...)
		default:
			http.Error(w, "Only GET and PUT are allowed", http.StatusMethodNotAllowed)
		}
	}
	if r.URL.Path == "/extends" {
		http.Error(w, "Not Implemented yet", http.StatusForbidden)
	}
}

func (h *UserHandler) handleGet(w http.ResponseWriter, r *http.Request, ids ...int) {
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
