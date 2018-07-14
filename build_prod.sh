#!/bin/bash
# docker run --rm -it -v "$GOPATH":/gopath -v "$(pwd)":/app -e "GOPATH=/gopath" -w /app golang:1.10 sh -c 'CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" -o abylebotter'
go build -a -installsuffix cgo -o abylebotter .
docker build -t abyle/projects/abylebotter-prod -f Dockerfile-prod .
