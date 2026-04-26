package local_test

import (
	"bytes"
	"io"
	"testing"

	localstorage "spec-streaming/backend/internal/storage/local"
)

func TestLocalStorageSavesAndReadsArtifacts(t *testing.T) {
	store := localstorage.New(t.TempDir())

	if err := store.SaveArtifact("videos/1/manifest.mpd", bytes.NewBufferString("manifest")); err != nil {
		t.Fatalf("save artifact: %v", err)
	}

	rc, err := store.OpenArtifact("videos/1/manifest.mpd")
	if err != nil {
		t.Fatalf("open artifact: %v", err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("read artifact: %v", err)
	}

	if string(data) != "manifest" {
		t.Fatalf("unexpected data: %s", string(data))
	}
}
