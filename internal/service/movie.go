package service

import (
	"context"
	"time"

	"watchlist-app/dto"
	"watchlist-app/ent"
	"watchlist-app/ent/movie"
	"watchlist-app/pkg/errors"
)

type MovieService struct {
	client *ent.Client
}

func NewMovieService(client *ent.Client) *MovieService {
	return &MovieService{
		client: client,
	}
}

// 映画リスト取得（フィルタリング付き）
func (s *MovieService) GetMovies(ctx context.Context, genreFilter, statusFilter, mediaTypeFilter string) ([]*ent.Movie, error) {
	query := s.client.Movie.Query()

	// ジャンルフィルタ
	if genreFilter != "" {
		query = query.Where(movie.GenreEQ(genreFilter))
	}

	// ステータスフィルター
	if statusFilter != "" {
		status := movie.WatchStatus(statusFilter)
		query = query.Where(movie.WatchStatusEQ(status))
	}

	// メディアタイプフィルタ
	if mediaTypeFilter != "" {
		mediaType := movie.MediaType(mediaTypeFilter)
		query = query.Where(movie.MediaTypeEQ(mediaType))
	}

	// 作成日時の降順でソート
	movies, err := query.Order(ent.Desc(movie.FieldCreatedAt)).All(ctx)
	if err != nil {
		return nil, errors.NewInternalServerError("映画の取得に失敗しました")
	}
	return movies, nil
}

// 映画詳細取得
func (s *MovieService) GetMovie(ctx context.Context, id int) (*ent.Movie, error) {
	movie, err := s.client.Movie.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NewNotFoundError("映画が見つかりません")
		}
		return nil, errors.NewInternalServerError("映画の取得に失敗しました")
	}
	return movie, nil
}

// 映画作成
func (s *MovieService) CreateMovie(ctx context.Context, req *dto.CreateMovieRequest) (*ent.Movie, error) {
	builder := s.client.Movie.Create().
		SetTitle(req.Title).
		SetMediaType(movie.MediaType(req.MediaType))

	// オプションフィールドの設定
	if req.Description != "" {
		builder = builder.SetDescription(req.Description)
	}

	if req.Genre != "" {
		builder = builder.SetGenre(req.Genre)
	}
	if req.ReleaseYear > 0 {
		builder = builder.SetReleaseYear(req.ReleaseYear)
	}
	if req.PosterURL != "" {
		builder = builder.SetPosterURL(req.PosterURL)
	}

	movie, err := builder.Save(ctx)
	if err != nil {
		return nil, errors.NewInternalServerError("映画の作成に失敗しました")
	}
	return movie, nil
}

// 映画更新
func (s *MovieService) UpdateMovie(ctx context.Context, id int, req *dto.UpdateMovieRequest) (*ent.Movie, error) {
	builder := s.client.Movie.UpdateOneID(id)

	// 更新するフィールドのみ設定
	if req.Title != "" {
		builder = builder.SetTitle(req.Title)
	}
	if req.Description != "" {
		builder = builder.SetDescription(req.Description)
	}
	if req.Genre != "" {
		builder = builder.SetGenre(req.Genre)
	}
	if req.ReleaseYear > 0 {
		builder = builder.SetReleaseYear(req.ReleaseYear)
	}
	if req.PosterURL != "" {
		builder = builder.SetPosterURL(req.PosterURL)
	}
	if req.MediaType != "" {
		builder = builder.SetMediaType(movie.MediaType(req.MediaType))
	}
	if req.WatchStatus != "" {
		builder = builder.SetWatchStatus(movie.WatchStatus(req.WatchStatus))

		// 視聴完了時は視聴完了日を設定
		if req.WatchStatus == "completed" {
			builder = builder.SetWatchedAt(time.Now())
		}
	}
	if req.Rating > 0 {
		builder = builder.SetRating(req.Rating)
	}
	if req.Review != "" {
		builder = builder.SetReview(req.Review)
	}

	movie, err := builder.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NewNotFoundError("映画が見つかりません")
		}
		return nil, errors.NewInternalServerError("映画の更新に失敗しました")
	}
	return movie, nil
}

// 映画削除
func (s *MovieService) DeleteMovie(ctx context.Context, id int) error {
	err := s.client.Movie.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.NewNotFoundError("映画が見つかりません")
		}
		return errors.NewInternalServerError("映画の削除に失敗しました")
	}
	return nil
}

// ジャンル一覧取得（統計用）
func (s *MovieService) GetGenres(ctx context.Context) ([]string, error) {
	var genres []string
	err := s.client.Movie.Query().
		GroupBy(movie.FieldGenre).
		Scan(ctx, &genres)
	if err != nil {
		return nil, errors.NewInternalServerError("ジャンルの取得に失敗しました")
	}
	return genres, err
}

// 視聴統計取得
func (s *MovieService) GetWatchStats(ctx context.Context) (map[string]int, error) {
	stats := make(map[string]int)

	// ステータス別カウント
	statuses := []movie.WatchStatus{
		movie.WatchStatusWantToWatch,
		movie.WatchStatusWatching,
		movie.WatchStatusCompleted,
		movie.WatchStatusDropped,
	}

	for _, status := range statuses {
		count, err := s.client.Movie.Query().
			Where(movie.WatchStatusEQ(status)).
			Count(ctx)
		if err != nil {
			return nil, errors.NewInternalServerError("統計情報の取得に失敗しました")
		}
		stats[string(status)] = count
	}

	return stats, nil
}
