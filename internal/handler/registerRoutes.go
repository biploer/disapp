package handler

import (
	"burning-notes/internal/config"
	"burning-notes/internal/usecase"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Dependences struct {
	AssetsFS fs.FS
	Message  usecase.Message
	Config   config.Config
}

func RegisterRoutes(router chi.Router, deps Dependences) {
	messageHandler := messageHandler{
		message: deps.Message,
		domain:  deps.Config.Domain,
		scheme:  deps.Config.Scheme,
	}

	router.Handle("/assets/*", http.StripPrefix("/assets", http.FileServer(http.FS(deps.AssetsFS))))
	router.Get("/", handler(homeHandler{}.handleIndex))
	router.Get("/m/{uuid}", handler(messageHandler.showPreview))
	router.Post("/api/messages/take", handler(messageHandler.showMessage))
	router.Post("/api/messages", handler(messageHandler.createMessage))
}
