package main

import (
	"fmt"
	"log"
	"net/http"
)

func serveHTTP(keyData *[]byte, port int) {
	fmt.Printf("Starting on port %d:...\n", port)
	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(fmt.Sprintf("%s\n", *keyData)))
		})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
