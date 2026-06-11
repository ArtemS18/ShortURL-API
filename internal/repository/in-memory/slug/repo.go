package slug

import (
	"context"
	"sync"

	"github.com/ArtemS18/ShortURL-API/internal/entity"
	"github.com/ArtemS18/ShortURL-API/internal/usecase/dto"
)

type SlugRepo struct {
	mx    sync.Mutex
	slugs map[string]string
}

func NewInMemorySlugRepo() *SlugRepo {
	return &SlugRepo{
		slugs: make(map[string]string),
	}
}

func (r *SlugRepo) GetURL(ctx context.Context, slug string) (*entity.URL, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	url, ok := r.slugs[slug]
	if !ok {
		return nil, &entity.NotFoundError{Field: "slug"}
	}
	return &entity.URL{Value: url}, nil

}

func (r *SlugRepo) CreateSlug(ctx context.Context, e *dto.CreateSlugDB) error {
	r.mx.Lock()
	defer r.mx.Unlock()
	if _, ok := r.slugs[e.Slug]; ok {
		return &entity.AlredyExitError{Field: "slug"}
	}
	r.slugs[e.Slug] = e.URL
	return nil
}
