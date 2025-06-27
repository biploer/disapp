package handler

import (
	"burning-notes/internal/view/layout"
	"burning-notes/internal/view/preshow"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type messageHandler struct {
	msgs   messStorage
	scheme string
	domain string
}

func (m messageHandler) createMessage(w http.ResponseWriter, r *http.Request) error {
	msg := r.FormValue("msg")
	dur, _ := time.ParseDuration("5m")
	msgId := m.msgs.Add(msg, dur)
	url := m.scheme + "://" + m.domain + "/m/" + msgId.String()
	return layout.MessageCard("Ссылка на записку", url).Render(r.Context(), w)
}

func (m messageHandler) showPreview(w http.ResponseWriter, r *http.Request) error {
	const INCORRECT_URL = "Sorry, but your URL is incorrect"

	id := r.PathValue("uuid")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		w.Write([]byte(INCORRECT_URL))
		return err
	}

	isExist := m.msgs.Check(parsedId)
	if !isExist {
		_, err = w.Write([]byte(INCORRECT_URL))
		return err
	}
	return preshow.Index(id).Render(r.Context(), w)
}

func (m messageHandler) showMessage(w http.ResponseWriter, r *http.Request) error {
	id := r.FormValue("uuid")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	msg, err := m.msgs.Take(parsedId)
	if err != nil {
		return err
	}

	return layout.MessageCard("Записка", msg).Render(r.Context(), w)
}
