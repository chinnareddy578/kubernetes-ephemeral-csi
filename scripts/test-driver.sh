#!/bin/bash

set -e

# Build and deploy the CSI driver
echo "Building and deploying CSI driver..."
./scripts/build.sh
./scripts/deploy.sh

# Wait for the CSI driver to be ready
echo "Waiting for CSI driver to be ready..."
kubectl wait --for=condition=ready pod -l app=ephemeral-csi-controller -n kube-system --timeout=300s
kubectl wait --for=condition=ready pod -l app=ephemeral-csi-node -n kube-system --timeout=300s

# Create test pod
echo "Creating test pod..."
kubectl apply -f deploy/kubernetes/test-pod.yaml

# Wait for pod to be ready
echo "Waiting for test pod to be ready..."
kubectl wait --for=condition=ready pod/test-ephemeral --timeout=300s

# Check pod logs
echo "Checking pod logs..."
kubectl logs test-ephemeral

# Check volume status
echo "Checking volume status..."
kubectl describe pod test-ephemeral | grep -A 5 "Volumes:"

# Clean up
echo "Cleaning up..."
kubectl delete pod test-ephemeral 