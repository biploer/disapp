package handler

import (
	"disapp/internal/storage"
	"disapp/internal/view/layout"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type messageHandler struct {
	msgs   storage.Messages
	scheme string
	domain string
}

func (m messageHandler) handleMessageView(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("uuid")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
		// slog.Error("faild to pars uuid from URL", slog.Attr{
		// 	Key:   "error",
		// 	Value: slog.StringValue(err.Error()),
		// })
		// w.Write([]byte("Sorry, but your URL is incorrect"))
	}
	msg, err := m.msgs.Take(parsedId)
	if err != nil {
		return err
	}
	w.Write([]byte(msg))
	return nil
}

func (m messageHandler) createMessage(w http.ResponseWriter, r *http.Request) error {
	msg := r.FormValue("msg")
	dur, _ := time.ParseDuration("5m")
	msgId := m.msgs.Add(msg, dur)
	url := m.scheme + "://" + m.domain + "/m/" + msgId.String()
	return layout.MessageLink(url).Render(r.Context(), w)
}
