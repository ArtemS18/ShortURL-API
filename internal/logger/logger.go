package logger

import (
	"os"

	"github.com/ArtemS18/ShortURL-API/config"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func New(cfg *config.ProjectConfig) *logrus.Logger {
	godotenv.Load(".env")
	log := logrus.New()

	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{})
	log.SetLevel(logrus.InfoLevel)
	log.WithFields(logrus.Fields{
		"node-id": cfg.Server.NodeID,
	}).Info("logger initialized")

	return log
}
