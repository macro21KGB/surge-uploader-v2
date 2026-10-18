#!/bin/sh

set -e

mkdir -p builds

echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o builds/surge-uploader-linux-amd64 app.go

echo "Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -o builds/surge-uploader-darwin-amd64 app.go

echo "Building for macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build -o builds/surge-uploader-darwin-arm64 app.go

echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -o builds/surge-uploader-windows-amd64.exe app.go

echo "Builds completed successfully in builds/"
