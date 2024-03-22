#!/bin/bash


rm build/dispatch_server || true
rm build/dispatch_server.zip || true
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o build/dispatch_server ./main/
open build



