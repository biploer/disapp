package dto

import "time"

type CreateMessageInput struct {
	Body     string
	Duration time.Duration
}

type CreateMessageOutput struct {
	ID string
}
