package handler

import "net/http"

func GEt(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)

	w.Write([]byte("Hello"))
}

func POSt(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)

	Add(2)

	w.Write([]byte("Hello"))
}

func Add() {
	return
}
