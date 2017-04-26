package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/health", healthCheck)
	err := http.ListenAndServe("localhost:8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}
