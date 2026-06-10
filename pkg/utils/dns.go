package utils

import (
	"fmt"

	"github.com/ArtemS18/ShortURL-API/config"
)

func BuildDSN(cfg config.Postgres) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Database,
		cfg.SslMode,
	)
}
