package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bradfitz/http2"
)

const version_string = "SNU-SNU 0.2"

var (
	// The font that will be used
	// TODO: Make this configurable
	font = "<link href='http://fonts.googleapis.com/css?family=Lato:300' rel='stylesheet' type='text/css'>"

	// The CSS style that will be used for directory listings and when rendering markdown pages
	// TODO: Make this configurable
	style = "body { background-color: #f0f0f0; color: #0b0b0b; font-family: 'Lato', sans-serif; font-weight: 300; margin: 3.5em; font-size: 1.3em; } a { color: #4010010; font-family: courier; } a:hover { color: #801010; } a:active { color: yellow; } h1 { color: #101010; }"

	// List of filenames that should be displayed instead of a directory listing
	// TODO: Make this configurable
	indexFilenames = []string{"index.html", "index.md", "index.txt"}
)

func main() {
	flag.Parse()

	addr := ":3000"
	cert := "dummycert.pem"
	key := "dummykey.pem"

	// TODO: Use traditional args/flag handling.
	//       Add support for --help and --version.

	if len(flag.Args()) >= 1 {
		addr = flag.Args()[0]
	}
	if len(flag.Args()) >= 2 {
		cert = flag.Args()[1]
	}
	if len(flag.Args()) >= 3 {
		key = flag.Args()[2]
	}

	fmt.Println(version_string)
	fmt.Println()
	fmt.Println("HTTP/2 web server for static content")
	fmt.Println()
	fmt.Println("[arg 1], server addr\t", addr)
	fmt.Println("[arg 2], cert file\t", cert)
	fmt.Println("[arg 3], key file\t", key)
	fmt.Println()

	mux := http.NewServeMux()

	registerHandlers(mux, ".")

	s := &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Enable HTTP/2 support
	http2.ConfigureServer(s, nil)

	log.Println("Ready")

	if err := s.ListenAndServeTLS(cert, key); err != nil {
		fmt.Printf("Fail: %s\n", err)
	}
}
