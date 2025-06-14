package handler

import (
	"burning-notes/internal/config"
	"io/fs"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type messStorage interface {
	Add(text string, duration ...time.Duration) uuid.UUID
	Check(id uuid.UUID) bool
	Take(id uuid.UUID) (string, error)
}

type Dependences struct {
	AssetsFS fs.FS
	Msgs     messStorage
	Config   config.Config
}

type handlerFunc func(http.ResponseWriter, *http.Request) error

func RegisterRoutes(router chi.Router, deps Dependences) {
	messageHandler := messageHandler{
		msgs:   deps.Msgs,
		domain: deps.Config.Domain,
		scheme: deps.Config.Scheme,
	}

	router.Handle("/assets/*", http.StripPrefix("/assets", http.FileServer(http.FS(deps.AssetsFS))))
	router.Get("/", handler(homeHandler{}.handleIndex))
	router.Get("/m/{uuid}", handler(messageHandler.showPreview))
	router.Post("/api/messages/take", handler(messageHandler.showMessage))
	router.Post("/api/messages", handler(messageHandler.createMessage))
}

func handler(h handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			handleError(w, r, err)
		}
	}
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	_ = w
	_ = r
	slog.Error("error during request", slog.String("err", err.Error()))
}
