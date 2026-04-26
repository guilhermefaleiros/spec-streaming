package s3

import (
	"fmt"
	"io"
)

type Storage struct{}

func New() *Storage { return &Storage{} }

func (s *Storage) SaveSource(string, io.Reader) error            { return fmt.Errorf("not implemented") }
func (s *Storage) OpenSource(string) (io.ReadCloser, error)      { return nil, fmt.Errorf("not implemented") }
func (s *Storage) SaveArtifact(string, io.Reader) error          { return fmt.Errorf("not implemented") }
func (s *Storage) OpenArtifact(string) (io.ReadCloser, error)    { return nil, fmt.Errorf("not implemented") }
func (s *Storage) ArtifactExists(string) (bool, error)           { return false, fmt.Errorf("not implemented") }
