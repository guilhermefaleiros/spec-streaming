package transcoding

import "context"

type Service interface {
	Transcode(ctx context.Context, sourceKey string, videoID string) (string, error)
}
