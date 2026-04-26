package jobs

type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

type Job struct {
	ID           string
	VideoID      string
	Status       Status
	Attempts     int
	ErrorMessage string
}
