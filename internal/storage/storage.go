package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"task/internal/model"
)

type Storage interface {
	Load() (model.State, error)
	Save(model.State) error
}

type FileStorage struct {
	mu   sync.RWMutex
	path string
}

func NewFileStorage(path string) *FileStorage {
	return &FileStorage{path: path}
}

func (s *FileStorage) Load() (model.State, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var st model.State

	if _, err := os.Stat(s.path); errors.Is(err, os.ErrNotExist) {
		st.NextID = 1
		st.Batches = make(map[int64]model.Batch)
		return st, nil
	}

	data, err := os.ReadFile(s.path)
	if err != nil {
		return st, err
	}

	if err := json.Unmarshal(data, &st); err != nil {
		return st, err
	}

	if st.Batches == nil {
		st.Batches = make(map[int64]model.Batch)
	}

	return st, nil
}

func (s *FileStorage) Save(st model.State) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}

	tmp := s.path + ".tmp"

	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}

	return os.Rename(tmp, s.path)
}
