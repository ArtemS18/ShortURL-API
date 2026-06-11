package slug

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/ArtemS18/ShortURL-API/internal/entity"
	"github.com/ArtemS18/ShortURL-API/internal/usecase/dto"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
)

func TestCreateSlug(t *testing.T) {
	t.Parallel()

	// Экранируем обновленный SQL-запрос
	query := regexp.QuoteMeta(`INSERT INTO slugs (id, slug, url) 
            VALUES ($1, $2, $3)
            ON CONFLICT (url) 
            DO UPDATE SET url = EXCLUDED.url
            RETURNING slug, (xmax = 0) AS is_inserted;`)

	ID := int64(123)
	input := dto.CreateSlugDB{
		Slug: "slug",
		ID:   ID,
		URL:  "http://er.ru",
	}

	tests := []struct {
		name      string
		input     dto.CreateSlugDB
		setupMock func(m pgxmock.PgxPoolIface)
		want      *dto.CreateSlugResponse
		wantErr   error
	}{
		{
			name:  "OK - New Slug Created",
			input: input,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"slug", "is_inserted"}).
					AddRow(input.Slug, true)
				m.ExpectQuery(query).
					WithArgs(input.ID, input.Slug, input.URL).
					WillReturnRows(rows)
			},
			want: &dto.CreateSlugResponse{
				SlugURL:   input.Slug,
				IsCreated: true,
			},
			wantErr: nil,
		},
		{
			name:  "OK - Slug Already Existed",
			input: input,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"slug", "is_inserted"}).
					AddRow("old-slug", false)

				m.ExpectQuery(query).
					WithArgs(input.ID, input.Slug, input.URL).
					WillReturnRows(rows)
			},
			want: &dto.CreateSlugResponse{
				SlugURL:   "old-slug",
				IsCreated: false,
			},
			wantErr: nil,
		},
		{
			name:  "Database Error",
			input: input,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(query).
					WithArgs(input.ID, input.Slug, input.URL).
					WillReturnError(errors.New("pg error"))
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repoService := NewSlugRepo(mock)
			res, err := repoService.CreateSlug(context.Background(), &test.input)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				require.Nil(t, res)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, res)
		})
	}
}

func TestGetURL(t *testing.T) {
	t.Parallel()
	query := regexp.QuoteMeta(`SELECT url FROM slugs WHERE slug=$1`)
	slugInput := "test-slug"

	expectedURL := entity.URL{
		Value: "http://er.ru",
	}

	tests := []struct {
		name      string
		slug      string
		setupMock func(m pgxmock.PgxPoolIface)
		want      *entity.URL
		wantErr   error
	}{
		{
			name: "OK",
			slug: slugInput,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"url"}).
					AddRow(expectedURL.Value)

				m.ExpectQuery(query).
					WithArgs(slugInput).
					WillReturnRows(rows)
			},
			want:    &expectedURL,
			wantErr: nil,
		},
		{
			name: "Not Found",
			slug: slugInput,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(query).
					WithArgs(slugInput).
					WillReturnError(errors.New("not found"))
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name: "Database Error",
			slug: slugInput,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(query).
					WithArgs(slugInput).
					WillReturnError(errors.New("connection timeout"))
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewSlugRepo(mock)

			res, err := repo.GetURL(context.Background(), &entity.Slug{Value: test.slug})
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				require.Nil(t, res)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, res)
		})
	}
}
