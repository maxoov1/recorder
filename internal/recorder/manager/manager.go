package manager

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/maxoov1/recorder/internal/recorder"
)

var (
	_defaultPlaylist = "rolling_stream.m3u8"
)

type RecorderManager struct {
	base string

	instancesMutex sync.Mutex
	instances      map[string]*recorder.Recorder
}

func New(base string) *RecorderManager {
	return &RecorderManager{
		base: filepath.Clean(base), instances: make(map[string]*recorder.Recorder),
	}
}

func (m *RecorderManager) Run(ctx context.Context, identifier, endpoint string) error {
	m.instancesMutex.Lock()
	defer m.instancesMutex.Unlock()

	if _, exist := m.instances[identifier]; exist {
		return fmt.Errorf("instance %q already exist", identifier)
	}

	baseIdentifier := filepath.Join(m.base, identifier)

	if _, err := os.Stat(baseIdentifier); errors.Is(err, os.ErrNotExist) {
		log.Printf("folder %q doesn't exist - creating...", baseIdentifier)

		if err := os.MkdirAll(baseIdentifier, 0o777); err != nil {
			return fmt.Errorf("failed to create folder %q: %w", baseIdentifier, err)
		}
	}

	instance := recorder.New(
		identifier, endpoint, filepath.Join(baseIdentifier, _defaultPlaylist))

	if err := instance.Run(ctx); err != nil {
		return fmt.Errorf("failed to run instance: %w", err)
	}

	m.instances[identifier] = instance

	return nil
}

func (m *RecorderManager) Stop(identifier string) error {
	m.instancesMutex.Lock()
	defer m.instancesMutex.Unlock()

	instance, exist := m.instances[identifier]
	if !exist {
		return fmt.Errorf("instance %q doesn't exist", identifier)
	}

	if err := instance.Stop(); err != nil {
		return fmt.Errorf("failed to stop instance: %w", err)
	}

	delete(m.instances, identifier)

	return nil
}
