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

// path can be ""
func serveFile(mux *http.ServeMux, fsroot, urlroot, urlpath, filename string, serveThese []string) {
	log.Println("Registering handler for " + urlpath)
	log.Println("  fsroot:", fsroot)
	log.Println("  urlroot:", urlroot)
	log.Println("  filename:", filename)
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
			registerIndexDir(mux, dirname, urlpath, serveThese)
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
	log.Println("Registering index handler for " + urlpath)
	log.Println("  For these files:", served)
	mux.HandleFunc(urlpath, func(w http.ResponseWriter, req *http.Request) {
		style := "body { background-color: darkgray; text: red; font-family: courier; }"
		fmt.Fprint(w, "<html><body style=\""+style+"\">")
		for _, filename := range served {
			fmt.Fprint(w, "<p><a href=\"/"+filename+"\">"+filename+"</a></p><br>")
		}
		fmt.Fprint(w, "<body></html>")
	})
}

func registerIndexDir(mux *http.ServeMux, fsroot, urlroot string, filetypes []string) {
	log.Println("Registering index dir:")
	log.Println("  fsroot:", fsroot)
	log.Println("  urlroot:", urlroot)
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
			served = append(served, filename)
			serveFile(mux, fsroot, urlroot, urlroot+filename, filename, filetypes)
			if (filename == "index.html") || (filename == "index.txt") {
				serveFile(mux, fsroot, urlroot, urlroot, filename, filetypes)
				//if (urlroot != "/") && strings.HasSuffix(urlroot, "/") {
				//	serveFile(mux, fsroot, urlroot[:len(urlroot)-1], "/", filename, filetypes)
				//}
				hasIndexHandler = true
			}
		}
	}
	if !hasIndexHandler {
		// If no index.html is found, create a nice webpage for showing all the files
		registerIndexPage(mux, urlroot, served)
		//if (urlroot != "/") && strings.HasSuffix(urlroot, "/") {
		//	registerIndexPage(mux, urlroot[:len(urlroot)-1], served)
		//}
	}
}

// Serve all files in the current directory, or only a few select filetypes (html, css, js, png and txt)
func registerHandlers(mux *http.ServeMux, allFiles bool) {
	var serveThese []string
	if !allFiles {
		serveThese = []string{"html", "css", "js", "png", "txt", "md"}
	}
	registerIndexDir(mux, ".", "/", serveThese)
}
