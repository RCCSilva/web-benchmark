package main

import (
	"encoding/json"
	"log"
	"net/http"
)

const jsonContentType = "application/json"

type response struct {
	Status string
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(response{Status: "ok"})
}

func main() {
	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(IndexHandler))
	log.Fatal(http.ListenAndServe(":3000", router))
}
