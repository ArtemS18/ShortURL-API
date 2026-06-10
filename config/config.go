package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var Config ProjectConfig

type (
	ProjectConfig struct {
		AppEnv       string `yaml:"app-env" env:"APP_ENV" env-default:"dev"`
		DatabaseType string `yaml:"database-type" env:"DATABASE_TYPE" env-default:"pgsql"`
		CORS         `yaml:"cors"`
		Server       `yaml:"server"`
		Postgres     `yaml:"postgres"`
	}

	CORS struct {
		AllowedOrigins []string `yaml:"allowed-origins"`
		AllowedMethods []string `yaml:"allowed-methods"`
		AllowedHeaders []string `yaml:"allowed-headers"`
	}

	Server struct {
		Host         string        `yaml:"host"    env:"SRV_HOST" env-default:"localhost"`
		Port         int           `yaml:"port"    env:"SRV_PORT" env-default:"8000"`
		WriteTimeout time.Duration `yaml:"write-timeout"    env:"SRV_WRITE_TM" env-default:"5s"`
		ReadTimeout  time.Duration `yaml:"read-timeout"    env:"SRV_READ_TM" env-default:"5s"`
		IdleTimeout  time.Duration `yaml:"idle-timeout"    env:"SRV_IDLE_TM" env-default:"20s"`
		NodeID       int64         `yaml:"node-id"    env:"SRV_NODE_ID" env-default:"1"`
		BaseURL      string        `yaml:"base-url"    env:"SRV_BASE_URL" env-default:"http://example.com"`
	}

	Postgres struct {
		Host            string        `yaml:"host"    env:"PG_HOST" env-default:"localhost"`
		Port            int           `yaml:"port" env:"PG_PORT" env-default:"5432"`
		User            string        `yaml:"user" env:"PG_USER" env-default:"thebugs"`
		Password        string        `yaml:"password" env:"PG_PASSWORD" env-default:"thebugs"`
		Database        string        `yaml:"database" env:"PG_DB" env-default:"main"`
		SslMode         string        `yaml:"sslmode" env:"PG_SSL_MODE" env-default:"disable"`
		MaxOpenConns    int           `yaml:"max-open-connections" env:"PG_MAX_OPEN_CONN" env-default:"10"`
		ConnMaxLifetime time.Duration `yaml:"conn-max-lifetime" env:"PG_CONN_MAX_LIFETIME" env-default:"30s"`
	}
)

func Read(log *logrus.Logger) error {
	var err error

	if err := godotenv.Load(".env"); err != nil {
		log.Warnf("No .env file found: %v", err)
	}
	if err = cleanenv.ReadConfig("config/config.yaml", &Config); err != nil {
		return fmt.Errorf("error while reading application configuration: %w", err)
	}

	if err = cleanenv.ReadEnv(&Config); err != nil {
		return fmt.Errorf("error creating configuration object: %w", err)
	}
	log.Print(Config.CORS.AllowedOrigins)
	log.Println("reading configuration is successful")
	return nil
}
