package handler

import (
	"disapp/internal/config"
	"disapp/internal/storage"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type handlerFunc func(http.ResponseWriter, *http.Request) error

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
		msg := msgs.Take(parsedId) // TODO: add an existence check
		w.Write([]byte(msg.Body))
	})
	router.HandleFunc("POST /api/messages", func(w http.ResponseWriter, r *http.Request) {
		msg := r.FormValue("msg")
		dur, _ := time.ParseDuration("5m")
		msgId := msgs.Add(msg, dur)
		url := "http://" + config.Address + "/m/" + msgId.String()
		w.Write([]byte("<div style=\"background-color: bisque; width: auto; height: 100px;\">" + url + "</div>"))
	})

	router.Get("/", handler(homeHandler{}.handleIndex))
}

func handler(h handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			handleError(w, r, err)
		}
	}
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("error during request", slog.String("err", err.Error()))
}
