#!/bin/bash
export CGO_ENABLED=0 # static build
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o replish-linux-amd64