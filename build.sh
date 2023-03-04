#!/bin/bash

# Go to the current directory
cd "$(dirname "$0")"

kill $(lsof -t -i:8080)

# Build the Go project
go build

# Run the Go project in the background
./"$(basename "$(pwd)")" &

# Change directories to the React app
cd "$REACT_APP_PATH"

# Build the React app
npm start build
