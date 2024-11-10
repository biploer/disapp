package storage

import (
	"time"

	"github.com/google/uuid"
)

const DefaultDuration = "30m"

type Message struct {
	Body      string
	Duration  time.Duration
	CreatedAt time.Time
}

type Messages map[uuid.UUID]Message

func New() Messages {
	return make(Messages)
}

func (m Messages) Add(text string, duration ...time.Duration) uuid.UUID {
	id, _ := uuid.NewRandom()
	currDuration, _ := time.ParseDuration(DefaultDuration)

	if len(duration) > 0 {
		currDuration = duration[0]
	}

	m[id] = Message{
		Body:      text,
		Duration:  currDuration,
		CreatedAt: time.Now(),
	}
	return id
}

func (m Messages) Take(id uuid.UUID) Message {
	msg := m[id]
	delete(m, id)
	return msg
}
