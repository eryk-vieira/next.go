package handler

import "net/http"

func GET(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)

	w.Write([]byte("Hello"))
}

func POST(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)

	Add(2)

	w.Write([]byte("Hello"))
}

func POST(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)

	Add(2)

	w.Write([]byte("Hello"))
}

func Add(int) {
	return
}
