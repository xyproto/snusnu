package main

import (
	"fmt"
	"net/http"
)

func index(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "hi")
}

func registerHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/", index)
}
