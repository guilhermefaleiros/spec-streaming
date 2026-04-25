package videos

import "fmt"

type Status string

const (
	StatusUploaded   Status = "uploaded"
	StatusProcessing Status = "processing"
	StatusReady      Status = "ready"
	StatusFailed     Status = "failed"
)

type Video struct {
	ID               string
	Title            string
	OriginalFilename string
	Status           Status
	SourceStorageKey string
	ManifestKey      string
	ErrorMessage     string
}

func (v *Video) MarkProcessing() error {
	if v.Status != StatusUploaded {
		return fmt.Errorf("cannot move %s to processing", v.Status)
	}
	v.Status = StatusProcessing
	return nil
}

func (v *Video) MarkReady(manifestKey string) error {
	if v.Status != StatusProcessing {
		return fmt.Errorf("cannot move %s to ready", v.Status)
	}
	v.Status = StatusReady
	v.ManifestKey = manifestKey
	v.ErrorMessage = ""
	return nil
}

func (v *Video) MarkFailed(message string) {
	v.Status = StatusFailed
	v.ErrorMessage = message
}
