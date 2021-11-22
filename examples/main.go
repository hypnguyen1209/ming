package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hypnguyen1209/ming"
)

func homeHandle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome!"))
}

func userHandle(w http.ResponseWriter, r *http.Request) {
	response := fmt.Sprintf("Hello, %s", ming.GetParams(r, "name"))
	w.Write([]byte(response))
}

func idHandle(w http.ResponseWriter, r *http.Request) {
	response := fmt.Sprintf("Current path: %s", ming.GetParams(r, "path"))
	w.Write([]byte(response))
}

func main() {
	router := ming.Create()

	router.Get("/home", homeHandle)
	router.Get("/user/{name}", userHandle)
	router.Get("/id/*path", idHandle)
	log.Fatal(http.ListenAndServe(":8080", router))
}
