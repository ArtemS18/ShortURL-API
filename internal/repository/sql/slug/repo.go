package slug

import (
	"context"
	"fmt"

	"github.com/ArtemS18/ShortURL-API/internal/entity"
	repo "github.com/ArtemS18/ShortURL-API/internal/repository/sql"
	"github.com/ArtemS18/ShortURL-API/internal/usecase"
	"github.com/ArtemS18/ShortURL-API/internal/usecase/dto"
	"github.com/jackc/pgx/v5"
)

type SlugRepo struct {
	pool repo.DB
}

func NewSlugRepo(pool repo.DB) usecase.SlugRepository {
	return &SlugRepo{
		pool: pool,
	}
}

func (r *SlugRepo) GetURL(ctx context.Context, slug *entity.Slug) (*entity.URL, error) {
	sql := `SELECT url FROM slugs WHERE slug=$1`
	row, err := r.pool.Query(ctx, sql, slug.Value)

	if err != nil {
		return nil, repo.HandelPgErrors(err, "slug")
	}

	defer row.Close()

	urlEntity, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[entity.URL])
	if err != nil {
		return nil, repo.HandelPgErrors(err, "slug")
	}

	return &urlEntity, nil

}

func (r *SlugRepo) CreateSlug(ctx context.Context, e *dto.CreateSlugDB) (*dto.CreateSlugResponse, error) {
	sql := `INSERT INTO slugs (id, slug, url) 
            VALUES ($1, $2, $3)
            ON CONFLICT (url) 
            DO UPDATE SET url = EXCLUDED.url
            RETURNING slug, (xmax = 0) AS is_inserted;`

	row, err := r.pool.Query(ctx, sql, e.ID, e.Slug, e.URL)
	if err != nil {
		return nil, repo.HandelPgErrors(err, "slug")
	}
	defer row.Close()

	if row.Next() {
		var returnedSlug string
		var isInserted bool
		err = row.Scan(&returnedSlug, &isInserted)
		if err != nil {
			return nil, repo.HandelPgErrors(err, "slug")
		}
		return &dto.CreateSlugResponse{
			SlugURL:   returnedSlug,
			IsCreated: isInserted,
		}, nil
	}
	return nil, fmt.Errorf("no rows returned from insert")
}
