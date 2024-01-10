package main

import (
	"fmt"
	"net/http"
)

func handFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello goblog")
}

func main() {
	http.HandleFunc("/", handFunc)
	http.ListenAndServe(":8082", nil)
}
