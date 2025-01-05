package app

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"

	"disapp/internal/config"
	"disapp/internal/handler"
	"disapp/internal/storage"

	"github.com/go-chi/chi/v5"
)

func Run(ctx context.Context) error {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	projectBaseDir, err := filepath.Abs(filepath.Join("cmd", "disapp", "main.go"))
	if err != nil {
		panic(err)
	}
	projectBaseDir = filepath.Dir(filepath.Dir(filepath.Dir(projectBaseDir)))
	slog.Debug(fmt.Sprintf("Project base dir: %s", projectBaseDir))

	defaultConfigPath := filepath.Join(projectBaseDir, "config", "local.yaml")

	configPath := flag.String("p", defaultConfigPath, "set config path")
	flag.Parse()

	config := config.MustLoad(*configPath)
	msgs := storage.New()

	router := chi.NewRouter()
	handler.RegisterRoutes(router, handler.Dependences{
		AssetsFS:       http.Dir(config.AssetsDir),
		Msgs:           msgs,
		Config:         config,
		ProjectBaseDir: projectBaseDir,
	})

	server := http.Server{
		Addr:         config.Address,
		Handler:      router,
		ReadTimeout:  config.Timeout,
		WriteTimeout: config.Timeout,
		IdleTimeout:  config.IdleTimeout,
	}

	go func() {
		<-ctx.Done()
		slog.Warn("Shutting down server")
		server.Shutdown(ctx)
	}()

	startingMsg := fmt.Sprintf("--- Starting server on %s", config.Address)
	slog.Info(startingMsg)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
