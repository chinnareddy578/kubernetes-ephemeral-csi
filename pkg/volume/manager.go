package volume

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"k8s.io/klog/v2"
)

const (
	// Default volume permissions
	defaultVolumePermissions = 0755
)

// VolumeManager handles the lifecycle of ephemeral volumes
type VolumeManager struct {
	baseDir string
	mu      sync.RWMutex
	volumes map[string]*Volume
}

// Volume represents an ephemeral volume
type Volume struct {
	ID         string
	Path       string
	Size       int64
	PodID      string
	Retention  string
	MountPoint string
	SubPath    string
	Usage      int64
	LastAccess int64
}

// NewVolumeManager creates a new volume manager
func NewVolumeManager(baseDir string) (*VolumeManager, error) {
	if err := os.MkdirAll(baseDir, defaultVolumePermissions); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %v", err)
	}

	return &VolumeManager{
		baseDir: baseDir,
		volumes: make(map[string]*Volume),
	}, nil
}

// CreateVolume creates a new ephemeral volume
func (m *VolumeManager) CreateVolume(req *csi.CreateVolumeRequest) (*Volume, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Generate unique volume ID
	volumeID := generateVolumeID(req.Name)

	// Create volume directory
	volumePath := filepath.Join(m.baseDir, volumeID)
	if err := os.MkdirAll(volumePath, defaultVolumePermissions); err != nil {
		return nil, fmt.Errorf("failed to create volume directory: %v", err)
	}

	// Parse volume attributes
	size := parseSize(req.CapacityRange.GetRequiredBytes())
	retention := req.Parameters["retentionPolicy"]
	podID := req.Parameters["podID"]

	volume := &Volume{
		ID:        volumeID,
		Path:      volumePath,
		Size:      size,
		PodID:     podID,
		Retention: retention,
	}

	m.volumes[volumeID] = volume
	klog.Infof("Created volume %s at %s", volumeID, volumePath)

	return volume, nil
}

// DeleteVolume deletes an ephemeral volume
func (m *VolumeManager) DeleteVolume(volumeID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	volume, exists := m.volumes[volumeID]
	if !exists {
		return fmt.Errorf("volume %s not found", volumeID)
	}

	// Remove volume directory
	if err := os.RemoveAll(volume.Path); err != nil {
		return fmt.Errorf("failed to delete volume directory: %v", err)
	}

	delete(m.volumes, volumeID)
	klog.Infof("Deleted volume %s", volumeID)

	return nil
}

// GetVolume returns a volume by ID
func (m *VolumeManager) GetVolume(volumeID string) (*Volume, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	volume, exists := m.volumes[volumeID]
	if !exists {
		return nil, fmt.Errorf("volume %s not found", volumeID)
	}

	return volume, nil
}

// ListVolumes returns all volumes
func (m *VolumeManager) ListVolumes() []*Volume {
	m.mu.RLock()
	defer m.mu.RUnlock()

	volumes := make([]*Volume, 0, len(m.volumes))
	for _, volume := range m.volumes {
		volumes = append(volumes, volume)
	}

	return volumes
}

// UpdateVolumeUsage updates the usage statistics for a volume
func (m *VolumeManager) UpdateVolumeUsage(volumeID string, usage int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	volume, exists := m.volumes[volumeID]
	if !exists {
		return fmt.Errorf("volume %s not found", volumeID)
	}

	volume.Usage = usage
	return nil
}

// Helper functions

func generateVolumeID(name string) string {
	// TODO: Implement proper volume ID generation
	return fmt.Sprintf("vol-%s", name)
}

func parseSize(size int64) int64 {
	if size <= 0 {
		return 1 << 30 // Default to 1GB
	}
	return size
}
