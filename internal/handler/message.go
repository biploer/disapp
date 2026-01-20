package handler

import (
	"burning-notes/internal/dto"
	"burning-notes/internal/usecase"
	"burning-notes/internal/view/layout"
	"burning-notes/internal/view/preshow"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type messageHandler struct {
	message usecase.Message
	scheme  string
	domain  string
}

func (m messageHandler) createMessage(w http.ResponseWriter, r *http.Request) error {
	msg := r.FormValue("msg")
	dur, err := time.ParseDuration("5m")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unsaported type of duration"))
		return fmt.Errorf("invalid value of duration while handle creating message: %v", err)
	}

	input := dto.CreateMessageInput{
		Body:     msg,
		Duration: dur,
	}

	output, err := m.message.CreateMessage(&input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error occured while saving your message"))
		return fmt.Errorf("handling create message: %v", err)
	}

	url := m.scheme + "://" + m.domain + "/m/" + output.ID
	return layout.MessageCard("Ссылка на записку", url).Render(r.Context(), w)
}

func (m messageHandler) showPreview(w http.ResponseWriter, r *http.Request) error {
	const INCORRECT_URL = "Sorry, but your URL is incorrect"

	id := r.PathValue("uuid")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(INCORRECT_URL))
		return fmt.Errorf("invalid value of id while handle showPreview: %v", err)
	}

	output, err := m.message.CheckMessage(&dto.CheckMessageInput{ID: parsedId})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error occured while finding your message"))
		return err
	} else if !output.IsExist {
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
