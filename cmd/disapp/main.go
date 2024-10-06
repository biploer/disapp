package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"disapp/internal/config"
)

func main() {
	projectBaseDir, err := filepath.Abs(filepath.Join("cmd", "disapp", "main.go"))
	if err != nil {
		panic(err)
	}
	projectBaseDir = filepath.Dir(filepath.Dir(filepath.Dir(projectBaseDir)))
	slog.Info(fmt.Sprintf("Project base dir: %s", projectBaseDir))

	defaultConfigPath := filepath.Join(projectBaseDir, "config", "local.yaml")

	configPath := flag.String("p", defaultConfigPath, "set config path")
	flag.Parse()

	config := config.MustLoad(*configPath)

	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(*r)
		w.Write([]byte("<h1>What`s that?</h1>"))
	})
	router.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hi fellow!"))
	})

	server := http.Server{
		Addr:         config.Address,
		Handler:      router,
		ReadTimeout:  config.Timeout,
		WriteTimeout: config.Timeout,
		IdleTimeout:  config.IdleTimeout,
	}

	startingMsg := fmt.Sprintf("--- Starting server on %s", config.Address)
	slog.Info(startingMsg)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("failed to start server", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})
		os.Exit(1)
	}
}
