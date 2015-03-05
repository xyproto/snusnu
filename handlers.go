package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/russross/blackfriday"
)

var handled []string

// Check if a string is part of a list
func has(sl []string, s string) bool {
	for _, e := range sl {
		if e == s {
			return true
		}
	}
	return false
}

// path can be ""
func serveFile(mux *http.ServeMux, fsroot, urlroot, urlpath, filename string, serveThese []string) {
	//log.Println("Registering file handler for " + urlpath)
	//log.Println("  fsroot:", fsroot)
	//log.Println("  urlroot:", urlroot)
	//log.Println("  filename:", filename)
	if has(handled, urlpath) {
		//log.Println("Already handled:", urlpath)
		return
	}
	handled = append(handled, urlpath)
	mux.HandleFunc(urlpath, func(w http.ResponseWriter, req *http.Request) {
		// Store the filename in the closure
		file, err := os.Open(fsroot + "/" + filename)
		defer file.Close()
		if err != nil {
			fmt.Fprintf(w, "Error opening %s/%s: %s", fsroot, filename, err)
			return
		}

		fi, err := file.Stat()
		if err != nil {
			fmt.Fprintf(w, "Error running stat on %s/%s: %s", fsroot, filename, err)
			return
		}
		if fi.IsDir() {
			// Serve the entire directory
			dirname := fsroot + "/" + filename
			urlpath := urlroot + filename + "/"
			serveDir(mux, dirname, urlpath, serveThese)
			return
		}

		// Mimetypes
		if strings.HasSuffix(filename, ".html") {
			w.Header().Add("Content-Type", "text/html; charset=utf-8")
		} else if strings.HasSuffix(filename, ".css") {
			w.Header().Add("Content-Type", "text/css")
		} else if strings.HasSuffix(filename, ".js") {
			w.Header().Add("Content-Type", "application/javascript")
		} else if strings.HasSuffix(filename, ".txt") {
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		} else if strings.HasSuffix(filename, ".png") {
			w.Header().Add("Content-Type", "image/png")
		} else if strings.HasSuffix(filename, ".md") {
			w.Header().Add("Content-Type", "text/html")
			b, err := ioutil.ReadAll(file)
			if err != nil {
				fmt.Fprintf(w, "Unable to read %s: %s", filename, err)
				return
			}
			// Write to the ResponseWriter, from the Markdown converter
			io.Copy(w, bytes.NewBuffer(blackfriday.MarkdownCommon(b)))
			return
		}
		// Write to the ResponseWriter, from the File
		io.Copy(w, file)
	})

}

func registerIndexPage(mux *http.ServeMux, urlpath string, served []string) {
	//log.Println("Registering generated index page for " + urlpath)
	//log.Println("  For these files:", served)
	if has(handled, urlpath) {
		//log.Println("Already handled:", urlpath)
		return
	}
	handled = append(handled, urlpath)
	mux.HandleFunc(urlpath, func(w http.ResponseWriter, req *http.Request) {
		style := "body { background-color: #101010; text: #c0c0c0; font-family: tahoma,verdana,arial; margin: 3em; font-size: 1.5em; } a { color: #500000; } a:hover { color: #a00000; } a:active { color: #a0a0a0; }"
		fmt.Fprint(w, "<!doctype html><html><head><title>SNU SNU</title><style>"+style+"</style><head><body><h1>SNU SNU</h1>")
		for _, filename := range served {
			url := urlpath + filename
			fmt.Fprint(w, "<a href=\""+url+"\">"+filename+"</a><br>")
		}
		fmt.Fprint(w, "<body></html>")
	})
}

// Check if a filename that can already be opened is a directory
func directory(full_filename string) bool {
	fs, _ := os.Stat(full_filename)
	return fs.IsDir()
}

func serveDirWhenNeeded(mux *http.ServeMux, urlroot, fsroot, filename string, filetypes []string) {
	urlpath := urlroot + filename + "/"
	//log.Println("Serving directory when needed:", urlpath, urlroot, fsroot, filename)
	if has(handled, urlpath) {
		//log.Println("Already handled:", urlpath)
		return
	}
	handled = append(handled, urlpath)
	mux.HandleFunc(urlpath, func(w http.ResponseWriter, req *http.Request) {
		// Serve the entire directory
		fulldirname := fsroot + "/" + filename
		urlpath := urlroot + filename + "/"
		serveDir(mux, fulldirname, urlpath, filetypes)
	})
}

func serveDir(mux *http.ServeMux, fsroot, urlroot string, filetypes []string) {
	//log.Println("Registering dir handler for " + urlroot)
	//log.Println("  fsroot:", fsroot)
	//log.Println("  urlroot:", urlroot)

	hasIndexHandler := false
	dir, err := os.Open(fsroot)
	if err != nil {
		log.Fatalf("Could not open directory: %s (%s)", fsroot, err)
		return
	}
	filenames, err := dir.Readdirnames(-1)
	if err != nil {
		log.Fatalf("Could not read filenames from directory: %s (%s)", fsroot, err)
	}
	served := []string{}
	for _, filename := range filenames {
		serve := false
		if len(filetypes) == 0 {
			serve = true
		} else {
			// Check if the given filetype/extension should be served
			for _, filetype := range filetypes {
				if strings.HasSuffix(filename, "."+filetype) {
					serve = true
					break
				}
			}
		}
		if serve {
			if directory(fsroot + "/" + filename) {
				served = append(served, filename+"/")
				go serveDirWhenNeeded(mux, urlroot, fsroot, filename, filetypes)
			} else {
				served = append(served, filename)
				go serveFile(mux, fsroot, urlroot, urlroot+filename, filename, filetypes)
				if urlroot == "/" {
					if (filename == "index.html") || (filename == "index.txt") {
						go serveFile(mux, fsroot, urlroot, urlroot, filename, filetypes)
						hasIndexHandler = true
					}
				}
			}
		}
	}
	if !hasIndexHandler {
		// If no index.html is found, create a nice webpage for showing all the files
		go registerIndexPage(mux, urlroot, served)
	}
}

// Serve all files in the current directory, or only a few select filetypes (html, css, js, png and txt)
func registerHandlers(mux *http.ServeMux, allFiles bool) {
	var serveThese []string
	if !allFiles {
		serveThese = []string{"html", "css", "js", "png", "txt", "md"}
	}
	serveDir(mux, ".", "/", serveThese)
}
