package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, t *http.Request) {
	fmt.Printf("Hello World")
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
