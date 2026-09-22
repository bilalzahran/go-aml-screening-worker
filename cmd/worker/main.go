package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"cdaq-event-worker/internal/client/typesafe"
	"cdaq-event-worker/internal/config"
	"cdaq-event-worker/internal/handler"
	mongorepo "cdaq-event-worker/internal/repository/mongo"
	"cdaq-event-worker/internal/service"
	"cdaq-event-worker/internal/transport/rabbitmq"
	"cdaq-event-worker/internal/usecase"
)

var Version = "dev"

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	logger := buildLogger(cfg.LogLevel)
	logger.Info("starting cdaq-event-worker", "version", Version)

	// MongoDB
	mongoClient, err := mongo.Connect(options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		logger.Error("failed to connect to MongoDB", "error", err)
		os.Exit(1)
	}

	if err := mongoClient.Ping(context.Background(), nil); err != nil {
		logger.Error("failed to ping MongoDB", "error", err)
		os.Exit(1)
	}
	logger.Info("connected to MongoDB", "uri", cfg.MongoDB.URI)

	// Wire dependencies
	eventRepo := mongorepo.NewEventRepository(mongoClient, cfg.MongoDB.Database, cfg.MongoDB.Collection, logger)
	eventService := service.NewEventService(eventRepo, logger)

	jobRepository := mongorepo.NewJobRepository(mongoClient, cfg.MongoDB.Database, cfg.MongoDB.JobCollection, logger)
	childJobRepository := mongorepo.NewChildJobRepository(mongoClient, cfg.MongoDB.Database, cfg.MongoDB.ChildJobCollection, logger)
	jobService := service.NewJobService(jobRepository, childJobRepository, logger)

	screeningWcResultRepository := mongorepo.NewScreeningWcResultRepository(mongoClient, cfg.MongoDB.Database, cfg.MongoDB.ScreeningWcResultCollection, logger)
	screeningWcResultService := service.NewScreeningWcResultService(screeningWcResultRepository, logger)

	typesafeClient := typesafe.NewClientWithOptions(cfg.AIModel.TypesafeAPIKey, logger, typesafe.WithBaseURL(cfg.AIModel.BaseURL))

	handlers := map[string]handler.EventTypeHandler{
		"wc.ogs": usecase.NewOgsUseCase(logger, jobService, screeningWcResultService, typesafeClient),
	}
	eventHandler := handler.NewEventHandler(handlers, eventService, logger)
	consumer := rabbitmq.NewConsumer(cfg.RabbitMQ, eventHandler, logger)

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := consumer.Start(ctx); err != nil {
			logger.Error("consumer error", "error", err)
			stop()
		}
	}()

	logger.Info("worker started, waiting for events")
	<-ctx.Done()
	logger.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	consumer.Shutdown(shutdownCtx)

	if err := mongoClient.Disconnect(shutdownCtx); err != nil {
		logger.Error("error disconnecting MongoDB", "error", err)
	}

	logger.Info("shutdown complete")
}

func buildLogger(level string) *slog.Logger {
	var logLevel slog.Level
	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
}
