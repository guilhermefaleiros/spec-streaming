package jobs

import "context"

type Repository interface {
	Create(context.Context, *Job) error
	ClaimPending(context.Context) (*Job, error)
	Update(context.Context, *Job) error
}
