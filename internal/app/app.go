package app

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"

	conf "disapp/config"
	"disapp/internal/config"
	"disapp/internal/handler"
	"disapp/internal/storage"
	"disapp/web"

	"github.com/go-chi/chi/v5"
)

func Run(ctx context.Context) error {
	var isProdEnv bool

	flag.BoolVar(&isProdEnv, "p", false, "Use production env configuration, must exist disapp/config/prod.yaml")
	flag.Parse()

	if !isProdEnv {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Debug log level is on")
	}

	confFS := conf.ConfigFS()
	config := config.MustLoad(confFS, isProdEnv)

	msgs := storage.New()

	router := chi.NewRouter()
	assetsFS, err := web.AssetsFS()
	if err != nil {
		return err
	}

	handler.RegisterRoutes(router, handler.Dependences{
		AssetsFS: assetsFS,
		Msgs:     msgs,
		Config:   config,
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
