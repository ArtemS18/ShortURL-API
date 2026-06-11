package slug

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/ArtemS18/ShortURL-API/internal/entity"
	"github.com/ArtemS18/ShortURL-API/internal/usecase/dto"
	"github.com/ArtemS18/ShortURL-API/internal/usecase/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateSlug(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	tests := []struct {
		name       string
		inp        dto.CreateSlugRequest
		wantErr    error
		setupMocks func(slugRepo *mocks.MockSlugRepository, slugGen *mocks.MockSlugGenerator)
	}{
		{
			name:    "valid URL",
			inp:     dto.CreateSlugRequest{URL: "https://www.example.com"},
			wantErr: nil,
			setupMocks: func(slugRepo *mocks.MockSlugRepository, slugGen *mocks.MockSlugGenerator) {
				slugGen.EXPECT().GenerateSlug(&entity.URL{Value: "https://www.example.com"}).Return(&dto.CreateSlugDB{Slug: "abc123"}, nil)
				slugRepo.EXPECT().CreateSlug(gomock.Any(), gomock.Any()).Return(&dto.CreateSlugResponse{SlugURL: "12"}, nil)
			},
		},
		{
			name:       "invalid URL",
			inp:        dto.CreateSlugRequest{URL: "not a valid URL"},
			wantErr:    entity.InvalidInput,
			setupMocks: nil,
		},
		{
			name:    "repository error",
			inp:     dto.CreateSlugRequest{URL: "https://www.example.com"},
			wantErr: entity.ServiceError,
			setupMocks: func(slugRepo *mocks.MockSlugRepository, slugGen *mocks.MockSlugGenerator) {
				slugGen.EXPECT().GenerateSlug(&entity.URL{Value: "https://www.example.com"}).Return(&dto.CreateSlugDB{Slug: "abc123"}, nil)
				slugRepo.EXPECT().CreateSlug(gomock.Any(), gomock.Any()).Return(nil, entity.ServiceError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockSlugRepository(ctrl)
			gen := mocks.NewMockSlugGenerator(ctrl)
			if tt.setupMocks != nil {
				tt.setupMocks(repo, gen)
			}
			uc := NewSlugUseCase(repo, gen)
			_, err := uc.CreateSlug(ctx, &tt.inp)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}

		})
	}
}

func TestGetSlug(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	tests := []struct {
		name       string
		inp        dto.GetURLRequest
		wantErr    error
		setupMocks func(slugRepo *mocks.MockSlugRepository, slugGen *mocks.MockSlugGenerator)
	}{
		{
			name:    "valid URL",
			inp:     dto.GetURLRequest{SlugURL: "2jDafv-I6u"},
			wantErr: nil,
			setupMocks: func(slugRepo *mocks.MockSlugRepository, slugGen *mocks.MockSlugGenerator) {
				slugRepo.EXPECT().GetURL(gomock.Any(), &entity.Slug{Value: "2jDafv-I6u"}).Return(&entity.URL{Value: "https://www.example.com"}, nil)
			},
		},
		{
			name:       "invalid slug",
			inp:        dto.GetURLRequest{SlugURL: "aaaaaaaaaaaaaaaaaaa"},
			wantErr:    entity.InvalidInput,
			setupMocks: nil,
		},
		{
			name:    "repository error",
			inp:     dto.GetURLRequest{SlugURL: "2jDafv-I6u"},
			wantErr: entity.ServiceError,
			setupMocks: func(slugRepo *mocks.MockSlugRepository, slugGen *mocks.MockSlugGenerator) {
				slugRepo.EXPECT().GetURL(gomock.Any(), &entity.Slug{Value: "2jDafv-I6u"}).Return(nil, entity.ServiceError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockSlugRepository(ctrl)
			gen := mocks.NewMockSlugGenerator(ctrl)
			if tt.setupMocks != nil {
				tt.setupMocks(repo, gen)
			}
			uc := NewSlugUseCase(repo, gen)
			_, err := uc.GetURL(ctx, &tt.inp)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}

		})
	}
}

func TestValidateURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		inp     string
		wantErr bool
	}{
		{
			name:    "valid http URL with nums",
			inp:     "https://example123.com",
			wantErr: false,
		},
		{
			name:    "valid http URL with chars",
			inp:     "https://example-12-123.com",
			wantErr: false,
		},
		{
			name:    "valid http URL with ru",
			inp:     "https://привет.рф",
			wantErr: false,
		},
		{
			name:    "valid short URL",
			inp:     "https://1.ru",
			wantErr: false,
		},
		{
			name:    "valid long URL",
			inp:     "https://lavka.yandex.ru/catalog/promo/category/mars_ice_cream_may_june2026?erid=nyi26TK8Sq2EHf7MZe6np6X4SLbusF8y",
			wantErr: false,
		},
		{
			name:    "valid http URL",
			inp:     "http://example.com",
			wantErr: false,
		},
		{
			name:    "valid www URL",
			inp:     "http://www.example.com",
			wantErr: false,
		},
		{
			name:    "valid domains URL",
			inp:     "http://d2.ddd.ddd.example.com",
			wantErr: false,
		},
		{
			name:    "invalid URL",
			inp:     "invalid_url.ru",
			wantErr: true,
		},
		{
			name:    "invalid protocol",
			inp:     "smb://example.com",
			wantErr: true,
		},
		{
			name:    "invalid long url",
			inp:     fmt.Sprintf("https://%s.ru", strings.Repeat("a", MaxURLLength+1)),
			wantErr: true,
		},
		{
			name:    "invalid empty url",
			inp:     "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			uc := NewSlugUseCase(nil, nil)
			err := uc.validateURL(tt.inp)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateSlug(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		inp     string
		wantErr bool
	}{
		{
			name:    "valid slug",
			inp:     "2jDafvQI6u",
			wantErr: false,
		},
		{
			name:    "valid slug with chars",
			inp:     "2jDafv-I6u",
			wantErr: false,
		},
		{
			name:    "invalid long slug",
			inp:     "2jDafv-I6u2jDafv-I6u2jDafv-I6u",
			wantErr: true,
		},
		{
			name:    "invalid slug chars",
			inp:     "2jDafv-I6я",
			wantErr: true,
		},
		{
			name:    "invalid empty slug",
			inp:     "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			uc := NewSlugUseCase(nil, nil)
			err := uc.validateSlug(tt.inp)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
