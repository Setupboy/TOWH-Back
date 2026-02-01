package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Setupboy/TOWH-Back/internal/config"
	databse "github.com/Setupboy/TOWH-Back/internal/database"
	"github.com/Setupboy/TOWH-Back/internal/interfaces"
	"github.com/Setupboy/TOWH-Back/internal/logger"
	"github.com/Setupboy/TOWH-Back/internal/repositories"
	"github.com/Setupboy/TOWH-Back/internal/server"
	"github.com/Setupboy/TOWH-Back/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	log := logger.New()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}
	db, err := databse.New(&cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	mainDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	defer func(mainDB *sql.DB) {
		err = mainDB.Close()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to close main database")
		}
	}(mainDB)

	ctx := context.Background()

	gin.SetMode(cfg.Server.GinMode)

	userRepo := repositories.NewUserRepository(db)

	authService := service.NewAuthService(
		cfg,
		userRepo,
	)
	userService := service.NewUserService(db)

	var uploadProvider interfaces.UploadProvider

	//if cfg.Upload.UploadProvider == "s3" {
	//	uploadProvider = providers.NewS3Provider(cfg)
	//} else {
	//	uploadProvider = providers.NewLocalUploaderProvider(cfg.Upload.Path)
	//}

	uploadService := service.NewUploadService(uploadProvider)
	srv := server.New(cfg, db, &log, authService, userService, uploadService)

	router := srv.SetupRouter()

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Info().Str("Port", cfg.Server.Port).Msg("starting http server")
		if err = httpServer.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("failed to start http server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("starting server")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err = httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("failed to shutdown http server")
	}

	log.Info().Msg("shutting down database")
}
