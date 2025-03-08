package storage

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
)

const DEFAULT_DURATION = "30m"

type message struct {
	timer     *time.Timer
	Body      string
	CreatedAt time.Time
}

type storage struct {
	messages map[uuid.UUID]message
	sync.Mutex
}

func New() storage {
	return storage{messages: make(map[uuid.UUID]message)}
}

func (s *storage) Add(text string, duration ...time.Duration) uuid.UUID {
	s.Lock()
	defer s.Unlock()

	id, _ := uuid.NewRandom()
	currDuration, _ := time.ParseDuration(DEFAULT_DURATION)

	if len(duration) > 0 {
		currDuration = duration[0]
	}

	s.messages[id] = message{
		timer: time.AfterFunc(currDuration, func() {
			s.Lock()
			defer s.Unlock()

			delete(s.messages, id)
			slog.Debug("Message has been delete by timeout")
		}),
		Body:      text,
		CreatedAt: time.Now(),
	}
	slog.Debug("New message created")

	return id
}

func (s *storage) Take(id uuid.UUID) (string, error) {
	s.Lock()
	defer s.Unlock()

	msg, ok := s.messages[id]
	if !ok {
		return "", fmt.Errorf("message with current id didn`t exist")
	}
	msg.timer.Stop()
	delete(s.messages, id)

	slog.Debug("Message has been read")
	return msg.Body, nil
}
