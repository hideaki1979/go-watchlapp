package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	App      AppConfig
}

type AppConfig struct {
	Environment string
}

type ServerConfig struct {
	Host string
	Port string
}

type DatabaseConfig struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     int
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/")

	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// 必須キーを環境変数へ明示バインド
	_ = viper.BindEnv("app.environment")
	_ = viper.BindEnv("server.host")
	_ = viper.BindEnv("server.port")
	_ = viper.BindEnv("database.user")
	_ = viper.BindEnv("database.password")
	_ = viper.BindEnv("database.name")
	_ = viper.BindEnv("database.host")
	_ = viper.BindEnv("database.port")

	// セーフデフォルト
	viper.SetDefault("app.environment", "development")
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", "8080")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 設定ファイルが見つからない場合は環境変数のみから読み込む
		} else {
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
