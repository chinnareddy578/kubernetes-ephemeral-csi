# Kubernetes Ephemeral CSI Driver

This project implements a Container Storage Interface (CSI) driver for Kubernetes that supports ephemeral volumes. Ephemeral volumes are temporary storage volumes that are created and destroyed with the pod lifecycle.

## Architecture

The CSI driver consists of two main components:

- **Controller Service**: Handles volume provisioning and deletion.
- **Node Service**: Manages volume mounting and unmounting on the node.

The driver uses a local filesystem-based approach, where volumes are created as directories on the host. For ephemeral volumes, the driver creates the volume directory on the node if it does not exist, ensuring that the volume is available for the pod.

## CSI and Ephemeral Volumes

CSI is a standard interface for container orchestration systems to expose arbitrary storage systems to their container workloads. Ephemeral volumes are volumes that are created and destroyed with the pod lifecycle, providing temporary storage for applications.

## Usage

To use the CSI driver, you need to:

1. Deploy the CSI driver to your Kubernetes cluster.
2. Create a pod that uses the CSI driver for ephemeral storage.

Example pod specification:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test-ephemeral-volume
spec:
  containers:
  - name: test-container
    image: busybox
    command: ["/bin/sh", "-c", "while true; do echo 'Hello from ephemeral volume' >> /data/test.txt; sleep 5; done"]
    volumeMounts:
    - name: ephemeral-volume
      mountPath: /data
  volumes:
  - name: ephemeral-volume
    csi:
      driver: ephemeral.csi.local
      volumeAttributes:
        size: "1Gi"
```

## Setup and Running

### Prerequisites

- Kubernetes cluster (e.g., minikube)
- Docker
- Go 1.21 or later

### Building the Driver

1. Clone the repository:
   ```bash
   git clone https://github.com/chinnareddy578/kubernetes-ephemeral-csi.git
   cd kubernetes-ephemeral-csi
   ```

2. Build the driver:
   ```bash
   make build
   ```

3. Build the Docker image:
   ```bash
   docker build -t ephemeral-csi:latest .
   ```

4. Load the image into minikube:
   ```bash
   minikube image load ephemeral-csi:latest
   ```

### Deploying the Driver

1. Apply the CSI driver deployment:
   ```bash
   kubectl apply -f deploy/kubernetes/csi-driver.yaml
   ```

2. Verify the deployment:
   ```bash
   kubectl get pods -n kube-system | grep ephemeral-csi
   ```

### Testing the Driver

1. Deploy the test pod:
   ```bash
   kubectl apply -f deploy/kubernetes/test-pod.yaml
   ```

2. Check the pod status:
   ```bash
   kubectl get pods | grep test-ephemeral-volume
   ```

3. Verify the volume is mounted and writable:
   ```bash
   kubectl exec test-ephemeral-volume -- cat /data/test.txt
   ```

## Submitting Issues and Change Proposals

We welcome contributions! If you encounter any issues or have suggestions for improvements, please submit them through the GitHub issue tracker. For change proposals, please create a pull request with a detailed description of the changes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.