package entity

type URLInfo struct {
	Slug string `db:"slug"`
	URL  string `db:"url"`
	ID   int64  `db:"id"`
}
type Slug struct {
	Value string `db:"slug"`
}

type URL struct {
	Value string `db:"url"`
}
