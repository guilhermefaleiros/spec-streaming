package videos

import "context"

type Repository interface {
	Create(context.Context, *Video) error
	List(context.Context) ([]Video, error)
	GetByID(context.Context, string) (*Video, error)
	Update(context.Context, *Video) error
}
