#!/bin/sh

env GOOS=linux GOARCH=arm GOARM=7 go build cmd/main.go
