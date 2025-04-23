#!/bin/bash


rm build/taskDispatch.zip || true

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o build/taskDispatch ./main.go

cd build/ && zip taskDispatch.zip taskDispatch

scp -P3371 taskDispatch.zip root@154.91.231.31:/root/task_dispatch/

