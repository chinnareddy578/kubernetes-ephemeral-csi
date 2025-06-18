package main

import (
	"flag"
	"net"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/grpc"
	"k8s.io/klog/v2"

	"github.com/chinnareddy578/kubernetes-ephemeral-csi/pkg/driver"
	"github.com/container-storage-interface/spec/lib/go/csi"
)

var (
	endpoint = flag.String("endpoint", "unix:///var/lib/kubelet/plugins/ephemeral.csi.local/csi.sock", "CSI endpoint")
	nodeID   = flag.String("nodeid", "", "Node ID")
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	// Create the driver
	d, err := driver.NewDriver(*nodeID, "/var/lib/ephemeral-csi")
	if err != nil {
		klog.Fatalf("Failed to create driver: %v", err)
	}

	// Create the gRPC server
	s := grpc.NewServer()

	// Register the CSI services
	csi.RegisterIdentityServer(s, d)
	csi.RegisterControllerServer(s, d)
	csi.RegisterNodeServer(s, d)

	// Create the socket directory
	socketDir := filepath.Dir(strings.TrimPrefix(*endpoint, "unix://"))
	if err := os.MkdirAll(socketDir, 0755); err != nil {
		klog.Fatalf("Failed to create socket directory: %v", err)
	}

	// Remove the socket if it exists
	if err := os.Remove(strings.TrimPrefix(*endpoint, "unix://")); err != nil && !os.IsNotExist(err) {
		klog.Fatalf("Failed to remove existing socket: %v", err)
	}

	// Create the listener
	lis, err := net.Listen("unix", strings.TrimPrefix(*endpoint, "unix://"))
	if err != nil {
		klog.Fatalf("Failed to listen: %v", err)
	}

	// Create the registration file
	regFile := filepath.Join(socketDir, "registration")
	if err := os.WriteFile(regFile, []byte(`{"driverName":"ephemeral.csi.local","endpoint":"unix:///var/lib/kubelet/plugins/ephemeral.csi.local/csi.sock"}`), 0644); err != nil {
		klog.Fatalf("Failed to create registration file: %v", err)
	}

	// Start the server
	klog.Infof("Starting CSI driver on %s", *endpoint)
	if err := s.Serve(lis); err != nil {
		klog.Fatalf("Failed to serve: %v", err)
	}
}
