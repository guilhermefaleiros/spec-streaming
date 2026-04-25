package local

import (
	"io"
	"os"
	"path/filepath"
)

type Storage struct {
	root string
}

func New(root string) *Storage {
	return &Storage{root: root}
}

func (s *Storage) SaveSource(key string, reader io.Reader) error { return s.writeFile(key, reader) }
func (s *Storage) OpenSource(key string) (io.ReadCloser, error) { return os.Open(filepath.Join(s.root, key)) }
func (s *Storage) SaveArtifact(key string, reader io.Reader) error { return s.writeFile(key, reader) }
func (s *Storage) OpenArtifact(key string) (io.ReadCloser, error) { return os.Open(filepath.Join(s.root, key)) }

func (s *Storage) ArtifactExists(key string) (bool, error) {
	_, err := os.Stat(filepath.Join(s.root, key))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (s *Storage) writeFile(key string, reader io.Reader) error {
	path := filepath.Join(s.root, key)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, reader)
	return err
}
