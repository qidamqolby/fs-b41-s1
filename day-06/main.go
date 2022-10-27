package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()

	route.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello World"))
	}).Methods("GET")

	route.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello about"))
	}).Methods("GET")

	fmt.Println(("Server running on port 5000"))
	http.ListenAndServe("localhost:5000", route)
}
