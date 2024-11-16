package handler

import (
	"disapp/internal/storage"
	"net/http"

	"github.com/google/uuid"
)

type messageHandler struct {
	msgs storage.Messages
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
