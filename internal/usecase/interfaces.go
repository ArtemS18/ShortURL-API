package usecase

import (
	"context"

	"github.com/ArtemS18/ShortURL-API/internal/entity"
	"github.com/ArtemS18/ShortURL-API/internal/usecase/dto"
)

type SlugGenerator interface {
	GenerateSlug(url string) (*dto.CreateSlug, error)
}

type SlugRepository interface {
	CreateSlug(ctx context.Context, e *dto.CreateSlug) (*entity.URLInfo, error)
	GetURL(ctx context.Context, slug string) (*entity.URL, error)
}
