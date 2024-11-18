package handler

import (
	"disapp/internal/storage"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type messageHandler struct {
	msgs   storage.Messages
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
	msg := m.msgs.Take(parsedId) // TODO: add an existence check
	w.Write([]byte(msg.Body))
	return nil
}

func (m messageHandler) createMessage(w http.ResponseWriter, r *http.Request) error {
	msg := r.FormValue("msg")
	dur, _ := time.ParseDuration("5m")
	msgId := m.msgs.Add(msg, dur)
	url := "http://" + m.domain + "/m/" + msgId.String()
	w.Write([]byte("<div style=\"background-color: bisque; width: auto; height: 100px;\">" + url + "</div>"))
	return nil
}
