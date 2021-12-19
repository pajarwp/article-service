package article

import "time"

type ArticleRequest struct {
	Author  string `json:"author" validate:"required"`
	Title   string `json:"title" validate:"required"`
	Body    string `json:"body" validate:"required"`
	Created time.Time
}

type ArticleResponse struct {
	Author  string    `json:"author"`
	Title   string    `json:"title"`
	Body    string    `json:"body"`
	Created time.Time `json:"created"`
}

type QueryParams struct {
	Author string `json:"author"`
	Query  string `json:"query"`
}
