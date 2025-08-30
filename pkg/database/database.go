package database

import (
	"context"
	"fmt"
	"log"
	"watchlist-app/ent"
	"watchlist-app/pkg/config"

	_ "entgo.io/ent/dialect"
	_ "github.com/lib/pq"
)

type Database struct {
	Client *ent.Client
}

func New(cfg *config.Config) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
	)

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to postgres: %w", err)
	}

	// 開発環境でのみデバッグログを有効化
	if cfg.App.Environment == "development" {
		client = client.Debug()
	}

	return &Database{
		Client: client,
	}, nil
}

func (d *Database) AutoMigrate(ctx context.Context) error {
	log.Println("Running database migrations...")
	return d.Client.Schema.Create(ctx)
}

func (d *Database) Close() error {
	return d.Client.Close()
}

// ヘルスチェック用
func (d *Database) Ping(ctx context.Context) error {
	// Entクライアント経由での接続確認
	_, err := d.Client.Movie.Query().Count(ctx)
	return err
}
