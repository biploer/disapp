package usecase

import (
	"burning-notes/internal/dto"
	"fmt"
)

func (m Message) TakeMessage(input *dto.TakeMessageInput) (dto.TakeMessageOutput, error) {
	var output dto.TakeMessageOutput

	msg, err := m.messages.PopMessage(input.ID)
	if err != nil {
		return output, fmt.Errorf("pop message from storage: %v", err)
	}

	return dto.TakeMessageOutput{Body: msg}, nil
}
