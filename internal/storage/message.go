package storage

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

const DEFAULT_DURATION = "30m"

type Message struct {
	timer     *time.Timer
	Body      string
	CreatedAt time.Time
}

type Messages map[uuid.UUID]Message

func New() Messages {
	return make(Messages)
}

func (m Messages) Add(text string, duration ...time.Duration) uuid.UUID {
	id, _ := uuid.NewRandom()
	currDuration, _ := time.ParseDuration(DEFAULT_DURATION)

	if len(duration) > 0 {
		currDuration = duration[0]
	}

	m[id] = Message{
		timer: time.AfterFunc(currDuration, func() {
			delete(m, id)
			slog.Debug("Message has been delete by timeout")
		}),
		Body:      text,
		CreatedAt: time.Now(),
	}
	slog.Debug("New message created")

	return id
}

func (m Messages) Take(id uuid.UUID) (string, error) {
	msg, ok := m[id]
	if !ok {
		return "", fmt.Errorf("Message with current id didn`t exist")
	}
	msg.timer.Stop()
	delete(m, id)

	slog.Debug("Message has been read")
	return msg.Body, nil
}
