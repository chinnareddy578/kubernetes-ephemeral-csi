package volume

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"k8s.io/klog/v2"
)

// NodeMounter handles volume mounting operations
type NodeMounter struct {
	volumeManager *VolumeManager
}

// NewNodeMounter creates a new node mounter
func NewNodeMounter(volumeManager *VolumeManager) *NodeMounter {
	return &NodeMounter{
		volumeManager: volumeManager,
	}
}

// NodePublishVolume mounts the volume to the target path
func (m *NodeMounter) NodePublishVolume(req *csi.NodePublishVolumeRequest) error {
	volumeID := req.GetVolumeId()
	targetPath := req.GetTargetPath()

	// Get volume information
	volume, err := m.volumeManager.GetVolume(volumeID)
	if err != nil {
		return fmt.Errorf("failed to get volume: %v", err)
	}

	// Create target directory if it doesn't exist
	if err := os.MkdirAll(targetPath, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %v", err)
	}

	// Handle subpath if specified
	if req.GetVolumeContext()["subPath"] != "" {
		subPath := req.GetVolumeContext()["subPath"]
		volumePath := filepath.Join(volume.Path, subPath)

		// Create subpath directory
		if err := os.MkdirAll(volumePath, 0755); err != nil {
			return fmt.Errorf("failed to create subpath directory: %v", err)
		}

		// Bind mount the subpath
		if err := bindMount(volumePath, targetPath); err != nil {
			return fmt.Errorf("failed to bind mount subpath: %v", err)
		}
	} else {
		// Bind mount the entire volume
		if err := bindMount(volume.Path, targetPath); err != nil {
			return fmt.Errorf("failed to bind mount volume: %v", err)
		}
	}

	// Update volume mount point
	volume.MountPoint = targetPath
	klog.Infof("Mounted volume %s to %s", volumeID, targetPath)

	return nil
}

// NodeUnpublishVolume unmounts the volume from the target path
func (m *NodeMounter) NodeUnpublishVolume(req *csi.NodeUnpublishVolumeRequest) error {
	volumeID := req.GetVolumeId()
	targetPath := req.GetTargetPath()

	// Get volume information
	volume, err := m.volumeManager.GetVolume(volumeID)
	if err != nil {
		return fmt.Errorf("failed to get volume: %v", err)
	}

	// Unmount the volume
	if err := syscall.Unmount(targetPath, 0); err != nil {
		return fmt.Errorf("failed to unmount volume: %v", err)
	}

	// Remove target directory
	if err := os.RemoveAll(targetPath); err != nil {
		return fmt.Errorf("failed to remove target directory: %v", err)
	}

	// Clear volume mount point
	volume.MountPoint = ""
	klog.Infof("Unmounted volume %s from %s", volumeID, targetPath)

	return nil
}

// NodeGetVolumeStats returns volume statistics
func (m *NodeMounter) NodeGetVolumeStats(volumeID string) (*csi.NodeGetVolumeStatsResponse, error) {
	volume, err := m.volumeManager.GetVolume(volumeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get volume: %v", err)
	}

	// Get filesystem statistics
	var stat syscall.Statfs_t
	if err := syscall.Statfs(volume.Path, &stat); err != nil {
		return nil, fmt.Errorf("failed to get volume stats: %v", err)
	}

	// Calculate available and used space
	blockSize := stat.Bsize
	availableBytes := int64(stat.Bavail * uint64(blockSize))
	totalBytes := int64(stat.Blocks * uint64(blockSize))
	usedBytes := totalBytes - availableBytes

	return &csi.NodeGetVolumeStatsResponse{
		Usage: []*csi.VolumeUsage{
			{
				Unit:      csi.VolumeUsage_BYTES,
				Available: availableBytes,
				Total:     totalBytes,
				Used:      usedBytes,
			},
		},
	}, nil
}

// bindMount performs a bind mount operation
func bindMount(source, target string) error {
	// Use the mount command for bind mounting
	cmd := exec.Command("mount", "--bind", source, target)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to bind mount: %v", err)
	}
	return nil
}
