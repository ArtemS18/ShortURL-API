package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ArtemS18/ShortURL-API/config"
	_ "github.com/ArtemS18/ShortURL-API/docs"
	"github.com/ArtemS18/ShortURL-API/internal/delivery/restapi/middleware"
	slugHandler "github.com/ArtemS18/ShortURL-API/internal/delivery/restapi/slug"
	slugInMemoryRepo "github.com/ArtemS18/ShortURL-API/internal/repository/in-memory/slug"
	slugRepo "github.com/ArtemS18/ShortURL-API/internal/repository/sql/slug"
	"github.com/ArtemS18/ShortURL-API/internal/usecase"
	slugGenerator "github.com/ArtemS18/ShortURL-API/internal/usecase/generator"
	slugUseCase "github.com/ArtemS18/ShortURL-API/internal/usecase/slug"
	"github.com/ArtemS18/ShortURL-API/pkg/showflake"
	"github.com/ArtemS18/ShortURL-API/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           ShortURL API
// @version         1.0
// @description     Created by Artemii in 2026
// @termsOfService  http://swagger.io/terms/

// @license.name  MIT

// @host      localhost:8000

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func Run(cfg *config.ProjectConfig, log *logrus.Logger) {
	var repo usecase.SlugRepository
	switch cfg.DatabaseType {
	case "postgres":
		log.Info("Using PostgreSQL database")
		dsn := utils.BuildDSN(cfg.Postgres)
		pool, err := pgxpool.New(context.Background(), dsn)
		if err != nil {
			log.Fatalf("cannot create pgx pool: %v", err)
		}
		repo = slugRepo.NewSlugRepo(pool)
		defer pool.Close()
	case "in-memory":
		log.Info("Using in-memory database")
		repo = slugInMemoryRepo.NewInMemorySlugRepo()
	default:
		log.Fatalf("unsupported database type: %s", cfg.DatabaseType)
	}

	fkCfg := showflake.SnowflakeConfig{
		Epoch:         cfg.Epoch,
		NodeID:        cfg.Server.NodeID,
		TimestampBits: cfg.ShowFlake.TimestampBits,
		NodeBits:      cfg.ShowFlake.NodeBits,
		SequenceBits:  cfg.ShowFlake.SequenceBits,
	}
	fk, err := showflake.NewSnowflake(fkCfg)
	if err != nil {
		log.Fatalf("cannot create snowflake: %v", err)
	}

	gen := slugGenerator.NewSlugGeneratorUseCase(fk)

	uc := slugUseCase.NewSlugUseCase(repo, gen)
	handler := slugHandler.NewSlugHandler(uc)

	r := mux.NewRouter()
	LoggingMiddleware := middleware.LoggingMiddleware(log)
	r.Use(LoggingMiddleware)

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}).Methods(http.MethodGet)

	api := r.PathPrefix("/").Subrouter()
	CORS := cors.New(cors.Options{
		AllowedOrigins:   config.Config.CORS.AllowedOrigins,
		AllowedMethods:   config.Config.AllowedMethods,
		AllowedHeaders:   config.Config.AllowedHeaders,
		AllowCredentials: true,
	})
	api.Use(CORS.Handler)

	// v1 := api.PathPrefix("/v1").Subrouter()
	api.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)
	api.HandleFunc("/slugs", handler.CreateSlugHandler).Methods(http.MethodPost)
	api.HandleFunc("/{slug}", handler.GetURLHandler).Methods(http.MethodGet)

	serverAddress := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	srv := &http.Server{
		Handler:      r,
		Addr:         serverAddress,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	log.Infof("start listen: %s", serverAddress)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Warn("shutting down")
	os.Exit(0)

}
