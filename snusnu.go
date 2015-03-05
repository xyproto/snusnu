package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bradfitz/http2"
)

const version_string = "SNU-SNU 0.1"

func main() {
	flag.Parse()

	cert := "dummycert.pem"
	key := "dummykey.pem"
	addr := ":3000"

	if len(flag.Args()) == 1 {
		addr = flag.Args()[0]
	} else if len(flag.Args()) >= 3 {
		addr = flag.Args()[0]
		cert = flag.Args()[1]
		key = flag.Args()[2]
	}

	fmt.Println(version_string)
	fmt.Println()
	fmt.Println("HTTP/2 web server for static content")
	fmt.Println()
	fmt.Println("[addr]\t\t", addr)
	fmt.Println("[cert file]\t", cert)
	fmt.Println("[key file]\t", key)
	fmt.Println()

	mux := http.NewServeMux()

	registerHandlers(mux, true)

	s := &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Enable HTTP/2
	http2.ConfigureServer(s, nil)

	log.Println("Ready")

	err := s.ListenAndServeTLS("dummycert.pem", "dummykey.pem")
	if err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}

}
