# Kubernetes Ephemeral CSI Driver

A Container Storage Interface (CSI) driver for managing local ephemeral storage in Kubernetes clusters. This driver provides dynamic provisioning and lifecycle management of local ephemeral volumes for pods.

## Features

- Dynamic provisioning of local ephemeral volumes
- Volume lifecycle management (create/delete)
- Volume usage limits and monitoring
- Volume retention policies
- Disk pressure handling with eviction support
- Subpath support
- Full pod volume feature support

## Architecture

The driver consists of the following components:

1. CSI Plugin Driver
   - Implements the CSI interface
   - Handles volume operations (create/delete/attach/detach)
   - Manages volume lifecycle

2. Volume Manager
   - Handles volume provisioning
   - Manages volume limits
   - Implements retention policies
   - Handles disk pressure scenarios

3. Node Controller
   - Manages node-specific operations
   - Handles volume attachment/detachment
   - Monitors local storage usage

## Building

```bash
make build
```

## Installation

```bash
kubectl apply -f deploy/kubernetes/
```

## Usage

To use the ephemeral storage in your pods:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: example-pod
spec:
  containers:
  - name: example
    image: nginx
    volumeMounts:
    - name: ephemeral-storage
      mountPath: /data
  volumes:
  - name: ephemeral-storage
    csi:
      driver: ephemeral.csi.local
      volumeAttributes:
        size: "1Gi"
        retentionPolicy: "delete"
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.