.PHONY: build clean deploy

# Variables
IMAGE_NAME ?= ephemeral-csi-driver
IMAGE_TAG ?= latest
REGISTRY ?= your-registry

# Build the CSI driver
build:
	go build -o bin/ephemeral-csi-driver ./cmd/csi-driver

# Build the container image
image:
	docker build -t $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG) .

# Push the container image
push:
	docker push $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

# Deploy to Kubernetes
deploy:
	kubectl apply -f deploy/kubernetes/

# Clean up build artifacts
clean:
	rm -rf bin/
	go clean

# Run tests
test:
	go test -v ./...

# Run linter
lint:
	golangci-lint run

# Generate Kubernetes manifests
manifests:
	@echo "Generating Kubernetes manifests..."
	@sed "s|image: ephemeral-csi-driver:latest|image: $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)|g" deploy/kubernetes/csi-driver.yaml > deploy/kubernetes/csi-driver-generated.yaml

# All-in-one build and deploy
all: build image push manifests deploy 