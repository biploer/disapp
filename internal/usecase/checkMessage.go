package usecase

import (
	"burning-notes/internal/dto"
)

func (m Message) CheckMessage(input *dto.CheckMessageInput) (dto.CheckMessageOutput, error) {
	return dto.CheckMessageOutput{
		IsExist: m.messages.CheckMessage(input.ID),
	}, nil
}
