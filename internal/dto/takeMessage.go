package dto

import "github.com/google/uuid"

type TakeMessageInput struct {
	ID uuid.UUID
}

type TakeMessageOutput struct {
	Body string
}
