# Build stage
FROM golang:1.21 as builder

WORKDIR /app

# Copy source code
COPY . .

# Build the application
RUN go build -o ephemeral-csi cmd/driver/main.go

# Final stage
FROM ubuntu:22.04

# Install required packages
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy the binary from builder
COPY --from=builder /app/ephemeral-csi /ephemeral-csi

# Set the entrypoint
ENTRYPOINT ["/ephemeral-csi"] 