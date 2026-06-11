package usecase

import (
	"context"

	"github.com/ArtemS18/ShortURL-API/internal/entity"
	"github.com/ArtemS18/ShortURL-API/internal/usecase/dto"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks/mocks_repo.go -package=mocks
type SlugGenerator interface {
	GenerateSlug(url *entity.URL) (*dto.CreateSlugDB, error)
}

type SlugRepository interface {
	CreateSlug(ctx context.Context, e *dto.CreateSlugDB) (*dto.CreateSlugResponse, error)
	GetURL(ctx context.Context, slug *entity.Slug) (*entity.URL, error)
}
