package usecase

import (
	"burning-notes/internal/dto"
	"burning-notes/internal/storage"
	"fmt"
	"time"
)

func (m Message) CreateMessage(input *dto.CreateMessageInput) (dto.CreateMessageOutput, error) {
	var output dto.CreateMessageOutput
	duration, err := time.ParseDuration(storage.DEFAULT_DURATION)
	if err != nil {
		return output, fmt.Errorf("create message: %v", err)
	}

	if input.Duration != 0 {
		duration = input.Duration
	}

	return dto.CreateMessageOutput{
		ID: string(m.messages.AddMessage(input.Body, duration).String()),
	}, nil
}
