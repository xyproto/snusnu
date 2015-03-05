package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/bradfitz/http2"
)

func main() {
	flag.Parse()
	fmt.Println(flag.Args())

	mux := http.NewServeMux()

	registerHandlers(mux, true)

	s := &http.Server{
		Addr:           ":3000",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Enable HTTP/2
	http2.ConfigureServer(s, nil)

	err := s.ListenAndServeTLS("dummycert.pem", "dummykey.pem")
	if err != nil {
		fmt.Printf("Server failed: ", err.Error())
	}

}
