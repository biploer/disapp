package handler

import (
	"disapp/internal/config"
	"disapp/internal/storage"
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func RegisterRoutes(router chi.Router, msgs storage.Messages, config config.Config, projectBaseDir string) {
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
}
