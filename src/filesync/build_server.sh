#!/bin/sh

export GOPATH=$HOME/go:$HOME/github/ghostboyzone/goplayground
echo "Set GOPATH:" $GOPATH
echo "Clean build output..."
go clean
echo "Clean build output, done"
echo "Build filesync_server.go..."
go build filesync_server.go
echo "Build filesync_server.go, done"