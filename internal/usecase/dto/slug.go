package dto

type CreateSlugDB struct {
	URL  string
	Slug string
	ID   int64
}

type CreateSlugRequest struct {
	URL string `json:"url" validate:"required,url"`
}

type CreateSlugResponse struct {
	SlugURL   string `json:"slug"`
	IsCreated bool   `json:"-"`
}

type GetURLRequest struct {
	SlugURL string `json:"slug" validate:"required"`
}

type GetURLResponse struct {
	URL string `json:"url"`
}
