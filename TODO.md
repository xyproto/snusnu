* Figure out how to tell if the server is actually serving HTTP/2
    * Here is one way:
        `nghttp -v http://somewebpagesomewhere.com/ | grep GOAWAY`

