#!/bin/bash

set -e

# Create namespace if it doesn't exist
kubectl create namespace kube-system --dry-run=client -o yaml | kubectl apply -f -

# Deploy the CSI driver
echo "Deploying CSI driver..."
kubectl apply -f deploy/kubernetes/csi-driver.yaml
kubectl apply -f deploy/kubernetes/csi-controller.yaml
kubectl apply -f deploy/kubernetes/csi-node.yaml 