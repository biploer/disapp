package app

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log/slog"
	"net/http"

	conf "burning-notes/config"
	"burning-notes/internal/config"
	"burning-notes/internal/handler"
	"burning-notes/internal/storage"
	"burning-notes/web"

	"github.com/go-chi/chi/v5"
)

func Run(ctx context.Context) error {
	var isProdEnv bool

	flag.BoolVar(&isProdEnv, "p", false, "Use production env configuration, must exist burning-notes/config/prod.yaml")
	flag.Parse()

	if !isProdEnv {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Debug log level is on")
	}

	confFS := conf.ConfigFS()
	config := config.MustLoad(confFS, isProdEnv)

	msgs := storage.New()

	var tlsConfig *tls.Config
	if isProdEnv {
		cert, err := tls.LoadX509KeyPair(config.Certificate.Cert, config.Certificate.Key)
		if err != nil {
			return err
		}
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}

	router := chi.NewRouter()
	assetsFS, err := web.AssetsFS()
	if err != nil {
		return err
	}

	handler.RegisterRoutes(router, handler.Dependences{
		AssetsFS: assetsFS,
		Msgs:     &msgs,
		Config:   config,
	})

	server := http.Server{
		Addr:         config.Address,
		Handler:      router,
		ReadTimeout:  config.Timeout,
		WriteTimeout: config.Timeout,
		IdleTimeout:  config.IdleTimeout,
	}
	if isProdEnv {
		server.TLSConfig = tlsConfig
	}

	go func() {
		<-ctx.Done()
		slog.Warn("Shutting down server")
		server.Shutdown(ctx)
	}()

	startingMsg := fmt.Sprintf("--- Starting server on %s", config.Address)
	slog.Info(startingMsg)
	if isProdEnv {
		if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			return err
		}
	} else {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
	}

	return nil
}
