#!/bin/bash
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o replish-linux-amd64