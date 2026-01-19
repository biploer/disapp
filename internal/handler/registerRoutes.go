package handler

import (
	"burning-notes/internal/config"
	"io/fs"
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
