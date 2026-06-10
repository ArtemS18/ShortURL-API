package generator

import (
	"fmt"

	"github.com/ArtemS18/ShortURL-API/internal/usecase/dto"
	"github.com/ArtemS18/ShortURL-API/pkg/showflake"
)

const (
	SlugLength = 10
)

type SlugGeneratorUseCase struct {
	s *showflake.Snowflake
}

func NewSlugGeneratorUseCase(s *showflake.Snowflake) *SlugGeneratorUseCase {
	return &SlugGeneratorUseCase{
		s: s,
	}
}

func (uc *SlugGeneratorUseCase) GenerateSlug(url string) (*dto.CreateSlug, error) {
	id, err := uc.s.NextID()
	if err != nil {
		return nil, fmt.Errorf("uc.s.NextID: %w", err)
	}
	slug := uc.s.Int64ToBase63(id, SlugLength)
	return &dto.CreateSlug{
		URL:  url,
		Slug: slug,
		ID:   id,
	}, nil
}
