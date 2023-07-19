package main

import (
	"fmt"
	"net/http"
)

func handleAll(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "I wanna be slim!\n")
}

func main() {
	http.HandleFunc("/", handleAll)

	fmt.Printf("Listening on port 8090...")
	http.ListenAndServe(":8090", nil)
}
