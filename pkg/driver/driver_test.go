package driver

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDriver(t *testing.T) (*Driver, string) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "csi-test-*")
	require.NoError(t, err)

	driver, err := NewDriver("test-node-id", tempDir)
	require.NoError(t, err)
	require.NotNil(t, driver)

	return driver, tempDir
}

func cleanupTestDriver(t *testing.T, tempDir string) {
	err := os.RemoveAll(tempDir)
	require.NoError(t, err)
}

func TestNewDriver(t *testing.T) {
	driver, tempDir := setupTestDriver(t)
	defer cleanupTestDriver(t, tempDir)

	assert.Equal(t, driverName, driver.name)
	assert.Equal(t, driverVersion, driver.version)
	assert.Equal(t, "test-node-id", driver.nodeID)
	assert.Equal(t, tempDir, driver.basePath)
}

func TestGetPluginInfo(t *testing.T) {
	driver, tempDir := setupTestDriver(t)
	defer cleanupTestDriver(t, tempDir)

	req := &csi.GetPluginInfoRequest{}
	resp, err := driver.GetPluginInfo(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, driverName, resp.Name)
	assert.Equal(t, driverVersion, resp.VendorVersion)
}

func TestGetPluginCapabilities(t *testing.T) {
	driver, tempDir := setupTestDriver(t)
	defer cleanupTestDriver(t, tempDir)

	req := &csi.GetPluginCapabilitiesRequest{}
	resp, err := driver.GetPluginCapabilities(context.Background(), req)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Capabilities)
}

func TestProbe(t *testing.T) {
	driver, tempDir := setupTestDriver(t)
	defer cleanupTestDriver(t, tempDir)

	req := &csi.ProbeRequest{}
	resp, err := driver.Probe(context.Background(), req)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestNodeGetInfo(t *testing.T) {
	driver, tempDir := setupTestDriver(t)
	defer cleanupTestDriver(t, tempDir)

	req := &csi.NodeGetInfoRequest{}
	resp, err := driver.NodeGetInfo(context.Background(), req)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test-node-id", resp.NodeId)
}

func TestNodeGetCapabilities(t *testing.T) {
	driver, tempDir := setupTestDriver(t)
	defer cleanupTestDriver(t, tempDir)

	req := &csi.NodeGetCapabilitiesRequest{}
	resp, err := driver.NodeGetCapabilities(context.Background(), req)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Capabilities)
}

func TestCreateVolume(t *testing.T) {
	driver, tempDir := setupTestDriver(t)
	defer cleanupTestDriver(t, tempDir)

	req := &csi.CreateVolumeRequest{
		Name: "test-volume",
		CapacityRange: &csi.CapacityRange{
			RequiredBytes: 1024 * 1024 * 1024, // 1GB
		},
	}

	resp, err := driver.CreateVolume(context.Background(), req)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test-volume", resp.Volume.VolumeId)
	assert.Equal(t, int64(1024*1024*1024), resp.Volume.CapacityBytes)

	// Verify volume directory was created
	volumePath := filepath.Join(tempDir, "test-volume")
	_, err = os.Stat(volumePath)
	require.NoError(t, err)
}

func TestDeleteVolume(t *testing.T) {
	driver, tempDir := setupTestDriver(t)
	defer cleanupTestDriver(t, tempDir)

	// Create a test volume first
	volumePath := filepath.Join(tempDir, "test-volume")
	err := os.MkdirAll(volumePath, 0755)
	require.NoError(t, err)

	req := &csi.DeleteVolumeRequest{
		VolumeId: "test-volume",
	}

	resp, err := driver.DeleteVolume(context.Background(), req)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	// Verify volume directory was deleted
	_, err = os.Stat(volumePath)
	assert.True(t, os.IsNotExist(err))
}

func TestNodePublishVolume(t *testing.T) {
	// Skip test if not running as root
	if os.Geteuid() != 0 {
		t.Skip("Skipping test that requires root privileges")
	}

	driver, tempDir := setupTestDriver(t)
	defer cleanupTestDriver(t, tempDir)

	// Create a test volume first
	volumePath := filepath.Join(tempDir, "test-volume")
	err := os.MkdirAll(volumePath, 0755)
	require.NoError(t, err)

	targetPath := filepath.Join(tempDir, "target")
	err = os.MkdirAll(targetPath, 0755)
	require.NoError(t, err)

	req := &csi.NodePublishVolumeRequest{
		VolumeId:   "test-volume",
		TargetPath: targetPath,
		VolumeCapability: &csi.VolumeCapability{
			AccessType: &csi.VolumeCapability_Mount{
				Mount: &csi.VolumeCapability_MountVolume{},
			},
			AccessMode: &csi.VolumeCapability_AccessMode{
				Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
			},
		},
	}

	resp, err := driver.NodePublishVolume(context.Background(), req)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestNodeUnpublishVolume(t *testing.T) {
	// Skip test if not running as root
	if os.Geteuid() != 0 {
		t.Skip("Skipping test that requires root privileges")
	}

	driver, tempDir := setupTestDriver(t)
	defer cleanupTestDriver(t, tempDir)

	targetPath := filepath.Join(tempDir, "target")
	err := os.MkdirAll(targetPath, 0755)
	require.NoError(t, err)

	req := &csi.NodeUnpublishVolumeRequest{
		VolumeId:   "test-volume",
		TargetPath: targetPath,
	}

	resp, err := driver.NodeUnpublishVolume(context.Background(), req)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}
