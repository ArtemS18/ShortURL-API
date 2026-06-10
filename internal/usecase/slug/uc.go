package slug

import (
	"context"
	"fmt"

	"github.com/ArtemS18/ShortURL-API/config"
	"github.com/ArtemS18/ShortURL-API/internal/usecase"
	"github.com/ArtemS18/ShortURL-API/internal/usecase/dto"
)

type SlugUseCase struct {
	repo usecase.SlugRepository
	gen  usecase.SlugGenerator
}

func createSlugURL(slug string) string {
	return fmt.Sprintf("%s/%s", config.Config.Server.BaseURL, slug)
}

func NewSlugUseCase(repo usecase.SlugRepository, gen usecase.SlugGenerator) *SlugUseCase {
	return &SlugUseCase{
		repo: repo,
		gen:  gen,
	}
}

func (uc *SlugUseCase) GetURL(ctx context.Context, e *dto.GetURLRequest) (*dto.GetURLResponse, error) {
	urlEntity, err := uc.repo.GetURL(ctx, e.SlugURL)
	if err != nil {
		return nil, fmt.Errorf("uc.repo.GetURL: %w", err)
	}
	return &dto.GetURLResponse{URL: urlEntity.Value}, nil
}

func (uc *SlugUseCase) CreateSlug(ctx context.Context, e *dto.CreateSlugRequest) (*dto.CreateSlugResponse, error) {
	slugInfo, err := uc.gen.GenerateSlug(e.URL)
	if err != nil {
		return nil, fmt.Errorf("uc.gen.GenerateSlug: %w", err)
	}
	_, err = uc.repo.CreateSlug(ctx, slugInfo)
	if err != nil {
		return nil, fmt.Errorf("uc.repo.CreateSlug: %w", err)
	}
	slugURL := createSlugURL(slugInfo.Slug)
	return &dto.CreateSlugResponse{SlugURL: slugURL}, nil
}
