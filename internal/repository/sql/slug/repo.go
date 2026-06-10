package slug

import (
	"context"

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

func (r *SlugRepo) GetURL(ctx context.Context, slug string) (*entity.URL, error) {
	sql := `SELECT url FROM slugs WHERE slug=$1`
	row, err := r.pool.Query(ctx, sql, slug)

	if err != nil {
		return nil, repo.HandelPgErrors(err)
	}

	defer row.Close()

	urlEntity, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[entity.URL])
	if err != nil {
		return nil, repo.HandelPgErrors(err)
	}

	return &urlEntity, nil

}

func (r *SlugRepo) CreateSlug(ctx context.Context, e *dto.CreateSlug) (*entity.URLInfo, error) {
	sql := `INSERT INTO slugs (id, slug, url) VALUES ($1, $2, $3) RETURNING id, slug, url`
	row, err := r.pool.Query(ctx, sql, e.ID, e.Slug, e.URL)
	if err != nil {
		return nil, repo.HandelPgErrors(err)
	}
	defer row.Close()

	slugEntity, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[entity.URLInfo])
	if err != nil {
		return nil, repo.HandelPgErrors(err)
	}
	return &slugEntity, nil
}
