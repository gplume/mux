package main

import (
	"net/http"
	"reflect"
)

// HomeHandler ...
type HomeHandler struct{}

// HomeHandler implements ServeHTPP to return an http.Handler (interface satisfying)
func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logger.Println("|---------WELCOME TO THE API!---------|")

	q1 := r.URL.Query().Get("id")
	logger.Println("ID:", q1, reflect.TypeOf(q1))

	q2 := r.URL.Query().Get("page")
	logger.Println("PAGE:", q2, reflect.TypeOf(q2))

	// time.Sleep(2 * time.Second)
	// panic(errors.New("| XxXxXxX P A N I C XxXxXxX |"))

	JSON(w, http.StatusOK, JSMAP{
		"msg": "OK",
	})
}
