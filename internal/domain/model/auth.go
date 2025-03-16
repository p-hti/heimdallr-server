package model

import "time"

type User struct {
	ID       int64
	Email    string
	PassHash []byte
}

type LoginStruct struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshToken struct {
	ID        int64
	UserID    int64
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type ContextKey string
