# snusnu [![Build Status](https://travis-ci.org/xyproto/snusnu.svg?branch=master)](https://travis-ci.org/xyproto/snusnu)

Simple HTTP/2 server, for serving the files in a given directory.

**DEPRECATED!** I would rather recommend [Algernon](https://github.com/xyproto/algernon).


Features and limitations
------------------------

* Supports HTTP/2 and HTTPS.
* Supports Markdown when displaying files ending with `.md`.
* If `index.html`, `index.md` or `index.txt` is found, it will be used for the main page.
* It's reasonably fast. Runs as a native executable.
* Uses UTF-8 whenever possible.
* Only sets Content-Type for a few commonly used filetypes.
* Self-signed TLS certificates will make the browser complain, unless the certificates are imported somehow.

Usage
-----

`snusnu [directory] [host:port] [certfile] [keyfile]`

`host:port` can be just `:port` for localhost.

Examples
------------------------------

Share the current directory as https://localhost:3000/

`snusnu . :3000`

---

Share a single file as the main page at https://localhost/. This will make snusnu listen to port 443, which may need more permissions. Run as root or a dedicated user (recommended).

`./snusnu README.md`

Screenshot
----------

Old screenshot. The new design is cleaner.

<img src="https://raw.githubusercontent.com/xyproto/snusnu/master/img/snusnu.png">

General information
-------------------

* Version: 0.3
* License: MIT
* Alexander F Rødseth

