#!/bin/bash

set -e

# Build the CSI driver binary
echo "Building CSI driver..."
go build -o ephemeral-csi cmd/driver/main.go

# Build the container image
echo "Building container image..."
docker build -t ephemeral-csi:latest .

# Load the image into Minikube
echo "Loading image into Minikube..."
minikube image load ephemeral-csi:latest 