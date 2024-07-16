package main

import (
	"fmt"
	"net/http"

	service "github.com/eryk-vieira/workspace/src/services"
	"github.com/google/uuid"
)

func PUT(w http.ResponseWriter, r *http.Request) {
	fmt.Println(uuid.NewString())

	service.CalculateSomething()

	w.Write([]byte("Hello"))
}
