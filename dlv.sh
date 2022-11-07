#!/bin/sh

cd /usr/src
dlv debug --headless --listen=:2346 --api-version=2 --log ./cmd/http.go
