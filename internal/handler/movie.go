package handler

import (
	"net/http"
	"strconv"
	"watchlist-app/dto"
	"watchlist-app/ent"
	"watchlist-app/internal/service"
	"watchlist-app/pkg/errors"

	"github.com/labstack/echo/v4"
)

type MovieHandler struct {
	movieService *service.MovieService
}

func NewMovieHandler(movieService *service.MovieService) *MovieHandler {
	return &MovieHandler{
		movieService: movieService,
	}
}

// GET /api/v1/movies - 映画リスト取得
func (h *MovieHandler) GetMovies(c echo.Context) error {
	var filter dto.MovieFilter
	if err := c.Bind(&filter); err != nil {
		return errors.NewBadRequestError("クエリパラメータが正しくありません")
	}

	movies, err := h.movieService.GetMovies(c.Request().Context(), filter.Genre, filter.Status, filter.MediaType)
	if err != nil {
		return err
	}

	// レスポンス変換
	response := make([]*dto.MovieResponse, len(movies))
	for i, movie := range movies {
		response[i] = h.convertToMovieResponse(movie)
	}

	return c.JSON(http.StatusOK, dto.MoviesResponse{
		Data:  response,
		Count: len(response),
	})
}

// GET /api/v1/movies/:id - 映画詳細取得
func (h *MovieHandler) GetMovie(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("無効なIDです")
	}

	movie, err := h.movieService.GetMovie(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.MovieDetailResponse{
		Data: h.convertToMovieResponse(movie),
	})
}

// POST /api/v1/movies - 映画作成
func (h *MovieHandler) CreateMovie(c echo.Context) error {
	var req dto.CreateMovieRequest
	if err := c.Bind(&req); err != nil {
		return errors.NewBadRequestError("リクエストの形式が正しくありません")
	}

	if err := c.Validate(&req); err != nil {
		return errors.NewBadRequestError("入力値が正しくありません: " + err.Error())
	}

	movie, err := h.movieService.CreateMovie(c.Request().Context(), &req)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, dto.MovieDetailResponse{
		Data: h.convertToMovieResponse(movie),
	})
}

// PUT /api/v1/movies/:id - 映画更新
func (h *MovieHandler) UpdateMovie(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("無効なIDです")
	}

	var req dto.UpdateMovieRequest
	if err := c.Bind(&req); err != nil {
		return errors.NewBadRequestError("リクエストの形式が正しくありません")
	}

	if err := c.Validate(&req); err != nil {
		return errors.NewBadRequestError("入力値が正しくありません: " + err.Error())
	}

	movie, err := h.movieService.UpdateMovie(c.Request().Context(), id, &req)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.MovieDetailResponse{
		Data: h.convertToMovieResponse(movie),
	})
}

// DELETE /api/v1/movies/:id - 映画削除
func (h *MovieHandler) DeleteMovie(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("無効なIDです")
	}

	err = h.movieService.DeleteMovie(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "映画が正常に削除されました",
	})
}

// GET /api/v1/stats/genres - ジャンル一覧取得
func (h *MovieHandler) GetGenres(c echo.Context) error {
	genres, err := h.movieService.GetGenres(c.Request().Context())
	if err != nil {
		return nil
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"data": genres})
}

// GET /api/v1/stats/watch - 視聴統計取得
func (h *MovieHandler) GetWatchStats(c echo.Context) error {
	stats, err := h.movieService.GetWatchStats(c.Request().Context())
	if err != nil {
		return nil
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"data": stats})
}

// Entエンティティ → DTOレスポンスへの変換
func (h *MovieHandler) convertToMovieResponse(movie *ent.Movie) *dto.MovieResponse {
	return &dto.MovieResponse{
		ID:          movie.ID,
		Title:       movie.Title,
		Description: movie.Description,
		Genre:       movie.Genre,
		ReleaseYear: movie.ReleaseYear,
		PosterURL:   movie.PosterURL,
		MediaType:   string(movie.MediaType),
		WatchStatus: string(movie.WatchStatus),
		Rating:      movie.Rating,
		Review:      movie.Review,
		WatchedAt:   movie.WatchedAt,
		CreatedAt:   movie.CreatedAt,
		UpdatedAt:   movie.UpdatedAt,
	}
}
