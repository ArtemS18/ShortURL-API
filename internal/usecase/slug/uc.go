package slug

import (
	"context"
	"fmt"

	"github.com/ArtemS18/ShortURL-API/config"
	"github.com/ArtemS18/ShortURL-API/internal/entity"
	"github.com/ArtemS18/ShortURL-API/internal/usecase"
	"github.com/ArtemS18/ShortURL-API/internal/usecase/dto"
	"github.com/ArtemS18/ShortURL-API/pkg/utils"
)

var MaxURLLength = 2048
var MaxSlugLength = 10

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
	if err := uc.validateSlug(e.SlugURL); err != nil {
		return nil, fmt.Errorf("uc.validateSlug: %w", err)
	}
	urlEntity, err := uc.repo.GetURL(ctx, e.SlugURL)
	if err != nil {
		return nil, fmt.Errorf("uc.repo.GetURL: %w", err)
	}
	return &dto.GetURLResponse{URL: urlEntity.Value}, nil
}

func (uc *SlugUseCase) CreateSlug(ctx context.Context, e *dto.CreateSlugRequest) (*dto.CreateSlugResponse, error) {
	if err := uc.validateURL(e.URL); err != nil {
		return nil, fmt.Errorf("uc.validateURL: %w", err)
	}
	slugInfo, err := uc.gen.GenerateSlug(e.URL)
	if err != nil {
		return nil, fmt.Errorf("uc.gen.GenerateSlug: %w", err)
	}
	err = uc.repo.CreateSlug(ctx, slugInfo)
	if err != nil {
		return nil, fmt.Errorf("uc.repo.CreateSlug: %w", err)
	}
	slugURL := createSlugURL(slugInfo.Slug)
	return &dto.CreateSlugResponse{SlugURL: slugURL}, nil
}

func (uc *SlugUseCase) validateURL(url string) error {
	if url == "" {
		return entity.NewValidationError("url", "Cant be empty")
	}
	if len([]rune(url)) > MaxURLLength {
		return entity.NewValidationError("url", fmt.Sprintf("Too long (max: %d)", MaxURLLength))
	}
	if !utils.IsValidURL(url) {
		return entity.NewValidationError("url", "Invalid format")
	}
	return nil
}

func (uc *SlugUseCase) validateSlug(slug string) error {
	if slug == "" {
		return entity.NewValidationError("slug", "Cant be empty")
	}
	if len([]rune(slug)) > MaxSlugLength {
		return entity.NewValidationError("slug", fmt.Sprintf("Too long (max: %d)", MaxSlugLength))
	}
	if !utils.IsValidSlug(slug) {
		return entity.NewValidationError("slug", "Invalid format")
	}
	return nil
}
