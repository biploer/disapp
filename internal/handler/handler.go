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

type Dependences struct {
	AssetsFS       http.FileSystem
	Msgs           storage.Messages
	Config         config.Config
	ProjectBaseDir string
}

type handlerFunc func(http.ResponseWriter, *http.Request) error

func RegisterRoutes(router chi.Router, deps Dependences) {
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
		msg := deps.Msgs.Take(parsedId) // TODO: add an existence check
		w.Write([]byte(msg.Body))
	})
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
	slog.Error("error during request", slog.String("err", err.Error()))
}
