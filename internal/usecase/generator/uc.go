package generator

import (
	"fmt"

	"github.com/ArtemS18/ShortURL-API/internal/entity"
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

func (uc *SlugGeneratorUseCase) GenerateSlug(e *entity.URL) (*dto.CreateSlugDB, error) {
	id, err := uc.s.NextID()
	if err != nil {
		return nil, fmt.Errorf("uc.s.NextID: %w", err)
	}
	slug := uc.s.Int64ToBase63(id, SlugLength)
	return &dto.CreateSlugDB{
		URL:  e.Value,
		Slug: slug,
		ID:   id,
	}, nil
}
