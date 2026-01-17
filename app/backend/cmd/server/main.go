package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gorkem/md-to-pdf/internal/config"
	"github.com/gorkem/md-to-pdf/internal/handler"
	"github.com/gorkem/md-to-pdf/internal/middleware"
	"github.com/gorkem/md-to-pdf/internal/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Initialize logger
	logger := initLogger()
	defer logger.Sync()

	// Load configuration
	cfg := config.Load()

	// Ensure directories exist
	if err := ensureDirectories(cfg); err != nil {
		logger.Fatal("failed to create directories", zap.Error(err))
	}

	// Initialize services
	markdownService := service.NewMarkdownService()
	converterService := service.NewConverterService(cfg, markdownService, logger)
	pdfExtractorService := service.NewPDFExtractorService(cfg, logger)

	// Initialize handler
	h := handler.NewHandler(markdownService, converterService, pdfExtractorService, logger)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		BodyLimit:    cfg.Server.BodyLimit,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		ErrorHandler: middleware.ErrorHandler(logger),
	})

	// Setup middleware
	middleware.Setup(app, logger)

	// Register routes
	h.RegisterRoutes(app)

	// Serve static files
	app.Static("/output", cfg.Storage.OutputDir)

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logger.Info("shutting down server...")
		if err := app.Shutdown(); err != nil {
			logger.Error("server shutdown error", zap.Error(err))
		}
	}()

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	logger.Info("starting server",
		zap.String("address", addr),
		zap.String("version", "1.0.0"),
	)

	if err := app.Listen(addr); err != nil {
		logger.Fatal("server failed", zap.Error(err))
	}
}

func initLogger() *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}

	return logger
}

func ensureDirectories(cfg *config.Config) error {
	dirs := []string{cfg.Storage.TempDir, cfg.Storage.OutputDir}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}
