#SNU SNU [![Build Status](https://travis-ci.org/xyproto/snusnu.svg?branch=master)](https://travis-ci.org/xyproto/snusnu) [![GoDoc](https://godoc.org/github.com/xyproto/snusnu?status.svg)](http://godoc.org/github.com/xyproto/snusnu)

Simple HTTP/2 server, for serving the files in the current directory.


Features and limitations
------------------------

* Uses HTTP/2. It's still a new protocol. There might be bugs.
* If `index.html` or `index.txt` is found, it will be used for the main page.
* Supports Markdown when displaying files ending with `.md`.
* HTTPS requires the use of certificates. This will make the browsers complain, unless the certificates are added to the browser.
* Only sets Content-Type for a few commonly used filetypes.
* Uses UTF-8 whenever possible.
* Should be pretty fast.

Known bugs
----------

* Serving files from subdirectories currently does not work as it should.

General information
-------------------

* Version: 0.1
* License: MIT
* Alexander F Rødseth

