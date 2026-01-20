package dto

import "github.com/google/uuid"

type CheckMessageInput struct {
	ID uuid.UUID
}

type CheckMessageOutput struct {
	IsExist bool
}
