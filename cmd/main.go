package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Orbit API 🚀")
}

func main() {
	http.HandleFunc("/", handler)

	port := "8080"
	log.Println("Server starting on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
