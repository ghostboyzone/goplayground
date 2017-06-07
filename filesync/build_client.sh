#!/bin/sh

export GOPATH=$HOME/go:$HOME/github/ghostboyzone/goplayground
echo "Set GOPATH:" $GOPATH
echo "Clean build output..."
go clean
echo "Clean build output, done"
echo "Build filesync_client.go..."
go build filesync_client.go
echo "Build filesync_client.go, done"