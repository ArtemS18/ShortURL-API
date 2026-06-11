package delivery

import (
	"context"

	"github.com/ArtemS18/ShortURL-API/internal/usecase/dto"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks/mocks_uc.go -package=mocks
type SlugUseCase interface {
	CreateSlug(ctx context.Context, e *dto.CreateSlugRequest) (*dto.CreateSlugResponse, error)
	GetURL(ctx context.Context, slug *dto.GetURLRequest) (*dto.GetURLResponse, error)
}
