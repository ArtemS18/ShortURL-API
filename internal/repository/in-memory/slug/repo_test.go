package slug

import (
	"context"
	"sync"
	"testing"

	"github.com/ArtemS18/ShortURL-API/internal/entity"
	"github.com/ArtemS18/ShortURL-API/internal/usecase/dto"
	"github.com/stretchr/testify/require"
)

func TestInMemorySlugRepo_CreateSlugDB(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		setup   func(r *SlugRepo)
		input   dto.CreateSlugDB
		wantErr error
	}{
		{
			name:  "OK - Successfully created",
			setup: func(r *SlugRepo) {},
			input: dto.CreateSlugDB{
				ID:   1,
				Slug: "yandex",
				URL:  "https://ya.ru",
			},
			wantErr: nil,
		},
		{
			name: "Error - Slug already exists",
			setup: func(r *SlugRepo) {
				r.slugs["yandex"] = "https://google.com"
			},
			input: dto.CreateSlugDB{
				ID:   2,
				Slug: "yandex",
				URL:  "https://ya.ru",
			},
			wantErr: &entity.AlredyExitError{Field: "slug"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			repo := NewInMemorySlugRepo()
			test.setup(repo)

			err := repo.CreateSlug(context.Background(), &test.input)
			if test.wantErr != nil {
				require.ErrorAs(t, err, &test.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestInMemorySlugRepo_GetURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		setup   func(r *SlugRepo)
		slug    string
		want    *entity.URL
		wantErr error
	}{
		{
			name: "OK - Found",
			setup: func(r *SlugRepo) {
				r.slugs["go"] = "https://go.dev"
			},
			slug:    "go",
			want:    &entity.URL{Value: "https://go.dev"},
			wantErr: nil,
		},
		{
			name:    "Error - Not Found",
			setup:   func(r *SlugRepo) {},
			slug:    "missing",
			want:    nil,
			wantErr: &entity.NotFoundError{Field: "slug"},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			repo := NewInMemorySlugRepo()
			test.setup(repo)

			res, err := repo.GetURL(context.Background(), test.slug)
			if test.wantErr != nil {
				require.ErrorAs(t, err, &test.wantErr)
				require.Nil(t, res)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, res)
		})
	}
}

func TestConcurencyCreateSlugRepo(t *testing.T) {
	t.Parallel()
	var wg sync.WaitGroup

	input := dto.CreateSlugDB{
		ID:   1,
		Slug: "yandex",
		URL:  "https://ya.ru",
	}

	want := &entity.URL{Value: "https://ya.ru"}

	repo := NewInMemorySlugRepo()
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = repo.CreateSlug(context.Background(), &input)
		}()
	}
	wg.Wait()
	resp, err := repo.GetURL(context.Background(), input.Slug)
	require.NoError(t, err)
	require.Equal(t, resp, want)
	require.NoError(t, err)

}

func TestConcurencyGetSlugRepo(t *testing.T) {
	t.Parallel()
	var wg sync.WaitGroup

	exits := dto.CreateSlugDB{
		ID:   1,
		Slug: "yandex",
		URL:  "https://ya.ru",
	}

	want := &entity.URL{Value: "https://ya.ru"}
	repo := NewInMemorySlugRepo()
	err := repo.CreateSlug(context.Background(), &exits)
	require.NoError(t, err)

	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := repo.GetURL(context.Background(), exits.Slug)
			require.NoError(t, err)
			require.Equal(t, resp, want)
			require.NoError(t, err)

		}()
	}
	wg.Wait()
}
