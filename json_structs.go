package main

import (
	"time"

	"github.com/google/uuid"
)

type reqS struct {
	Body string `json:"body"`
}

type eS struct {
	Error string `json:"error"`
}

type cleanedBody struct {
	Cleaned_Body string `json:"cleaned_body"`
}

type emailReq struct {
	Email string `json:"email"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}
