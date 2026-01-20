package usecase

import (
	"time"

	"github.com/google/uuid"
)

type Storage interface {
	AddMessage(text string, duration time.Duration) uuid.UUID
	CheckMessage(id uuid.UUID) bool
	PopMessage(id uuid.UUID) (string, error)
}

type Message struct {
	messages Storage
}

func NewMessage(storage Storage) *Message {
	return &Message{
		messages: storage,
	}
}
