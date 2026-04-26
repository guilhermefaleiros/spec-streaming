package videos_test

import (
	"testing"

	"spec-streaming/backend/internal/videos"
)

func TestVideoStatusTransitions(t *testing.T) {
	video := videos.Video{Status: videos.StatusUploaded}

	if err := video.MarkProcessing(); err != nil {
		t.Fatalf("expected uploaded -> processing to pass: %v", err)
	}

	if err := video.MarkReady("videos/123/manifest.mpd"); err != nil {
		t.Fatalf("expected processing -> ready to pass: %v", err)
	}
}

func TestVideoRejectsInvalidTransition(t *testing.T) {
	video := videos.Video{Status: videos.StatusUploaded}

	if err := video.MarkReady("videos/123/manifest.mpd"); err == nil {
		t.Fatal("expected uploaded -> ready to fail")
	}
}
