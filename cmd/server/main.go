package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"watchlist-app/internal/router"
	"watchlist-app/pkg/config"
	"watchlist-app/pkg/database"
	"watchlist-app/pkg/errors"
	"watchlist-app/pkg/validator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 設定読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// DB接続
	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// マイグレーション実行
	ctx := context.Background()
	if err := db.AutoMigrate(ctx); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Echo インスタンス作成
	e := echo.New()

	// カスタムエラーハンドラ設定
	e.HTTPErrorHandler = customHTTPErrorHandler

	// バリデータ設定
	e.Validator = validator.New()

	// ミドルウェア設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: strings.Split(os.Getenv("CORS_ALLOW_ORIGINS"), ","),
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Content-Type",
		},
	}))

	// ヘルスチェックエンドポイント
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// ルート設定
	router.SetupRoutes(e, db.Client)

	// サーバー開始
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s", serverAddr)
	log.Printf("API Base URL: http://%s/api/v1", serverAddr)

	// Graceful shutdown
	go func() {
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// シグナル待機
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown server: %v", err)
	}

	log.Println("Server stopped")
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "Internal Server Error"

	if he, ok := err.(*errors.AppError); ok {
		code = he.Code
		message = he.Message
	} else if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = fmt.Sprint(he.Message)
	}

	// 既にレスポンスがコミットされている場合は何もしない
	if c.Response().Committed {
		return
	}

	// 本番では 5xx の詳細は固定文言にする
	if code >= 500 && os.Getenv("APP_ENVIRONMENT") == "production" {
		message = "Internal Server Error"
	}

	if err := c.JSON(code, map[string]interface{}{"code": code, "message": message}); err != nil {
		c.Logger().Error("Error handling failed:", err)
	}
}
