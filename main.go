package main

import (
	"fmt"
	"net/http"
)

func handFunc(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		fmt.Fprint(w, "根目录")
	} else if r.URL.Path == "/about" {
		fmt.Fprint(w, "about 目录")
	} else {
		fmt.Fprint(w, "404")
	}
	// fmt.Fprintf(w, "hello goblog")
}

func main() {
	http.HandleFunc("/", handFunc)
	http.ListenAndServe(":8082", nil)
}
