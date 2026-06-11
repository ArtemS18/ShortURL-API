package slug

import (
	"context"
	"sync"

	"github.com/ArtemS18/ShortURL-API/internal/entity"
	"github.com/ArtemS18/ShortURL-API/internal/usecase/dto"
)

type InsertResult struct {
	Slug       string
	IsInserted bool
}

type SlugRepo struct {
	mx    sync.Mutex
	slugs map[string]string
	urls  map[string]string
}

func NewInMemorySlugRepo() *SlugRepo {
	return &SlugRepo{
		slugs: make(map[string]string),
		urls:  make(map[string]string),
	}
}

func (r *SlugRepo) GetURL(ctx context.Context, slug *entity.Slug) (*entity.URL, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	url, ok := r.slugs[slug.Value]
	if !ok {
		return nil, &entity.NotFoundError{Field: "slug"}
	}
	return &entity.URL{Value: url}, nil
}

func (r *SlugRepo) CreateSlug(ctx context.Context, e *dto.CreateSlugDB) (*dto.CreateSlugResponse, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	if existingSlug, ok := r.urls[e.URL]; ok {
		return &dto.CreateSlugResponse{
			SlugURL:   existingSlug,
			IsCreated: false,
		}, nil
	}

	if _, ok := r.slugs[e.Slug]; ok {
		return nil, &entity.AlredyExitError{Field: "slug"}
	}
	r.slugs[e.Slug] = e.URL
	r.urls[e.URL] = e.Slug

	return &dto.CreateSlugResponse{
		SlugURL:   e.Slug,
		IsCreated: true,
	}, nil
}
