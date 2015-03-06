#SNU SNU [![Build Status](https://travis-ci.org/xyproto/snusnu.svg?branch=master)](https://travis-ci.org/xyproto/snusnu)

Simple HTTP/2 server, for serving the files in a given directory.


Features and limitations
------------------------

* Supports HTTP/2.
* Supports Markdown when displaying files ending with `.md`.
* If `index.html`, `index.md` or `index.txt` is found, it will be used for the main page.
* HTTPS requires the use of certificates. This will make the browsers complain, unless the certificates are added to the browser.
* Only sets Content-Type for a few commonly used filetypes.
* Uses UTF-8 whenever possible.
* Reasonably fast.

Usage
-----

`snusnu [path] [addr] [cert] [key]`

addr can just be ":port"

Examples
------------------------------

Simple commandline invocation:

`snusnu . :3000`


<!--
Screenshot
----------

<img src="https://raw.githubusercontent.com/xyproto/snusnu/master/img/snusnu.png">
-->

General information
-------------------

* Version: 0.2
* License: MIT
* Alexander F Rødseth

