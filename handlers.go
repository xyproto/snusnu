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

// When serving a file. The file must exist. Must be given a full filename.
func filePage(w http.ResponseWriter, filename string) {
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
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(w, "Unable to read %s: %s", filename, err)
			return
		}
		markdownBody := string(blackfriday.MarkdownCommon(b))
		fmt.Fprint(w, markdownPage(filename, markdownBody))
		return
	}
	// Write to the ResponseWriter, from the File
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		fmt.Fprintf(w, "Can't open %s: %s", filename, err)
	}
	// Serve the file
	io.Copy(w, file)
	return
}

// Directory listing
func directoryListing(w http.ResponseWriter, dirname string) {
	var buf bytes.Buffer
	sep := string(os.PathSeparator)
	for _, filename := range getFilenames(dirname) {

		// Find the full name
		full_filename := dirname
		if !strings.HasSuffix(full_filename, sep) {
			full_filename += sep
		}
		full_filename += filename

		// Output different entries for files and directories
		buf.WriteString(easyLink(filename, full_filename, isDir(full_filename)))
	}
	title := dirname
	// Strip the leading "./"
	if strings.HasPrefix(title, "./") {
		title = title[2:]
	}
	// Use the application title for the main page
	//if title == "" {
	//	title = version_string
	//}
	if buf.Len() > 0 {
		fmt.Fprint(w, easyPage(title, buf.String()))
	} else {
		fmt.Fprint(w, easyPage(title, "Empty directory"))
	}
}

// When serving a directory. The directory must exist. Must be given a full filename.
func dirPage(w http.ResponseWriter, dirname string) {
	// Handle the serving of index files, if needed
	for _, indexfile := range indexFilenames {
		filename := dirname + string(os.PathSeparator) + indexfile
		if exists(filename) {
			filePage(w, filename)
			return
		}
	}
	// Serve a directory listing of no index file is found
	directoryListing(w, dirname)
}

// When a file is not found
func noPage(filename string) string {
	return easyPage("Not found", "File not found: "+filename)
}

// Serve all files in the current directory, or only a few select filetypes (html, css, js, png and txt)
func registerHandlers(mux *http.ServeMux, servedir string) {
	log.Println("About to serve dir:", servedir)
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		urlpath := req.URL.Path
		filename := url2filename(servedir, urlpath)
		if !exists(filename) {
			fmt.Fprint(w, noPage(filename))
			return
		}
		if isDir(filename) {
			dirPage(w, filename)
			return
		}
		filePage(w, filename)
	})
}
