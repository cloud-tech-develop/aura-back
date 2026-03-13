package enterprise

import "context"

type Repository interface {
	Create(ctx context.Context, e *Enterprise) error
	GetBySlug(ctx context.Context, slug string) (*Enterprise, error)
	GetBySubDomain(ctx context.Context, subDomain string) (*Enterprise, error)
	GetByEmail(ctx context.Context, email string) (*Enterprise, error)
	List(ctx context.Context) ([]Enterprise, error)
	Update(ctx context.Context, e *Enterprise) error
	Delete(ctx context.Context, id int64) error
}
