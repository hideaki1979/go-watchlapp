package router

import (
	"watchlist-app/ent"
	"watchlist-app/internal/handler"
	"watchlist-app/internal/service"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, client *ent.Client) {
	// サービス初期化
	movieService := service.NewMovieService(client)

	// ハンドル初期化
	movieHandle := handler.NewMovieHandler(movieService)

	// API v1グループ
	api := e.Group("/api/v1")

	// 映画関連ルート
	movies := api.Group("/movies")
	movies.GET("", movieHandle.GetMovies)
	movies.POST("", movieHandle.CreateMovie)
	movies.GET("/:id", movieHandle.GetMovie)
	movies.PUT("/:id", movieHandle.UpdateMovie)
	movies.DELETE("/:id", movieHandle.DeleteMovie)

	// 統計用エンドポイント
	stats := api.Group("/stats")
	stats.GET("/genres", movieHandle.GetGenres)

	stats.GET("/watch", movieHandle.GetWatchStats)
}