#!/bin/sh
# For generating SSL certs, for testing/development.
# Just press return at all the prompt, but enter "localhost" at Common Name.
openssl req -x509 -newkey rsa:4096 -keyout dummykey.pem -out dummycert.pem -days 9999 -nodes
