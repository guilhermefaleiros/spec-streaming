package storage

import "io"

type Storage interface {
	SaveSource(key string, reader io.Reader) error
	OpenSource(key string) (io.ReadCloser, error)
	SaveArtifact(key string, reader io.Reader) error
	OpenArtifact(key string) (io.ReadCloser, error)
	ArtifactExists(key string) (bool, error)
}
