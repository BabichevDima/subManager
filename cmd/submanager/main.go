package main

import (
	"context"
	"os"
	"runtime"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"net/http"

	"github.com/BabichevDima/subManager/internal/config"
	"github.com/BabichevDima/subManager/internal/db"
	"github.com/BabichevDima/subManager/internal/http/handlers"
	"github.com/BabichevDima/subManager/internal/http/middleware"

	router "github.com/BabichevDima/subManager/internal/http"
	"github.com/BabichevDima/subManager/internal/repository"
	"github.com/BabichevDima/subManager/internal/usecase"
	"github.com/BabichevDima/subManager/pkg/graceful"
	"github.com/BabichevDima/subManager/pkg/logger"

	_ "github.com/BabichevDima/subManager/internal/docs"
)

func main() {
	_ = godotenv.Load(".env")

	logger.Init()
	defer logger.L.Sync()

	cfgPath := os.Getenv("CONFIG_DEV_PATH")

	err := config.Init(cfgPath)
	if err != nil {
		logger.Fatal("Config init error:", zap.Error(err))
	}

	logger.Info("Config loaded",
		zap.String("db_host", config.Cfg.DB.Host),
		zap.Int("db_port", config.Cfg.DB.Port),
	)

	dbConn, err := db.InitPostgres(config.Cfg.DB)
	if err != nil {
		logger.Fatal("Failed to init db", zap.Error(err))
	}
	logger.Info("Successfully connected to the database")

	subscriptionRepo := repository.NewSubscriptionRepository(dbConn)
	subscriptionUsecase := usecase.NewUserUsecase(subscriptionRepo)
	subscriptionHandler := handlers.NewUserHandler(subscriptionUsecase)

	mux := http.NewServeMux()
	router.RegisterRoutes(mux, subscriptionHandler)
	handler := middleware.RequestLogger(logger.L, mux)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Info("Application starting",
		zap.String("version", "1.0.0"),
		zap.String("go_version", runtime.Version()),
	)

	go func() {
		logger.Info("HTTP server is listening",
			zap.String("address", "http://localhost"+server.Addr),
			zap.String("docs", "http://localhost"+server.Addr+"/swagger"),
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	graceful.GracefulShutdown(
		cancel,
		server,
		logger.L,
		5*time.Second,
		db.ShutdownDB(dbConn),
	)
}
