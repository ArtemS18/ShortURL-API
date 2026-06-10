package dto

type CreateSlug struct {
	URL  string
	Slug string
	ID   int64
}

type CreateSlugRequest struct {
	URL string `json:"url" validate:"required,url"`
}

type CreateSlugResponse struct {
	SlugURL string `json:"slug"`
}

type GetURLRequest struct {
	SlugURL string `json:"slug" validate:"required"`
}

type GetURLResponse struct {
	URL string `json:"url"`
}
