package main

import (
	"github.com/ArtemS18/ShortURL-API/config"
	"github.com/ArtemS18/ShortURL-API/internal/app"
	"github.com/ArtemS18/ShortURL-API/internal/logger"
	"github.com/sirupsen/logrus"
)

func main() {
	launchLogger := logrus.New()
	err := config.Read(launchLogger)
	if err != nil {
		launchLogger.Fatalf("Config error: %s", err)
	}
	appLogger := logger.New(&config.Config)

	app.Run(&config.Config, appLogger)

}
