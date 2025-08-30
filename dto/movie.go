package dto

import "time"

// 映画・ドラマリクエスト作成用の構造体
type CreateMovieRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Genre       string `json:"genre"`
	ReleaseYear int    `json:"release_year"`
	PosterURL   string `json:"poster_url"`
	MediaType   string `json:"media_type" validate:"oneof=movie tv_series documentary anime"`
}

// 映画更新リクエスト
type UpdateMovieRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Genre       string `json:"genre"`
	ReleaseYear int    `json:"release_year"`
	PosterURL   string `json:"poster_url"`
	MediaType   string `json:"media_type" validate:"omitempty,oneof=movie tv_series documentary anime"`
	WatchStatus string `json:"watch_status" validate:"omitempty,oneof=want_to_watch watching completed dropped"`
	Rating      int    `json:"rating" validate:"omitempty,min=1,max=5"`
	Review      string `json:"review"`
}

// 映画レスポンス
type MovieResponse struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Genre       string    `json:"genre,omitempty"`
	ReleaseYear int       `json:"release_year,omitempty"`
	PosterURL   string    `json:"poster_url,omitempty"`
	MediaType   string    `json:"media_type"`
	WatchStatus string    `json:"watch_status"`
	Rating      int       `json:"rating,omitempty"`
	Review      string    `json:"review,omitempty"`
	WatchedAt   time.Time `json:"watched_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// フィルターパラメータ
type MovieFilter struct {
	Genre     string `query:"genre"`
	Status    string `query:"status"`
	MediaType string `query:"media_type"`
}

// レスポンスラッパー
type MoviesResponse struct {
	Data  []*MovieResponse `json:"data"`
	Count int              `json:"count"`
}

type MovieDetailResponse struct {
	Data *MovieResponse `json:"data"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
