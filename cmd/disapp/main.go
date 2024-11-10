package main

import (
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"disapp/internal/config"
	"disapp/internal/storage"

	"github.com/google/uuid"
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
	msgs := storage.New()

	router := http.NewServeMux()
	router.HandleFunc("/hi/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		w.Write([]byte("Hi, " + name))
	})
	router.HandleFunc("GET /m/{uuid}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("uuid")
		parsedId, err := uuid.Parse(id)
		if err != nil {
			slog.Error("faild to pars uuid from URL", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			w.Write([]byte("Sorry, but your URL is incorrect"))
		}
		msg := msgs.Take(parsedId) // TODO: надо добавить проверку на наличие
		w.Write([]byte(msg.Body))
	})
	router.HandleFunc("POST /api/messages", func(w http.ResponseWriter, r *http.Request) {
		msg := r.FormValue("msg")
		dur, _ := time.ParseDuration("5m")
		msgId := msgs.Add(msg, dur)
		url := "http://" + config.Address + "/m/" + msgId.String()
		w.Write([]byte("<div style=\"background-color: bisque; width: auto; height: 100px;\">" + url + "</div>"))
	})
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		name := "index.html"
		path := filepath.Join(projectBaseDir, "web", name)
		htmlTemp, err := template.New(name).ParseFiles(path)
		if err != nil {
			slog.Error("faild to pars html template", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
		}
		htmlTemp.Execute(w, nil)
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
