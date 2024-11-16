package handler

import (
	"disapp/internal/config"
	"disapp/internal/storage"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Dependences struct {
	AssetsFS       http.FileSystem
	Msgs           storage.Messages
	Config         config.Config
	ProjectBaseDir string
}

type handlerFunc func(http.ResponseWriter, *http.Request) error

func RegisterRoutes(router chi.Router, deps Dependences) {
	router.Get("/m/{uuid}", handler(messageHandler{msgs: deps.Msgs}.handleMessageView))
	router.HandleFunc("POST /api/messages", func(w http.ResponseWriter, r *http.Request) {
		msg := r.FormValue("msg")
		dur, _ := time.ParseDuration("5m")
		msgId := deps.Msgs.Add(msg, dur)
		url := "http://" + deps.Config.Address + "/m/" + msgId.String()
		w.Write([]byte("<div style=\"background-color: bisque; width: auto; height: 100px;\">" + url + "</div>"))
	})

	router.Get("/", handler(homeHandler{}.handleIndex))
	router.Handle("/assets/*", http.StripPrefix("/assets", http.FileServer(deps.AssetsFS)))
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
