package model

import "time"

type User struct {
	ID        int64     `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Name      string    `json:"name" db:"name"`
	PassHash  []byte    `json:"pass_hash" db:"pass_hash"`
	CreatedAt time.Time `json:"created_at"`
}

type SignInStruct struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SingUpStruct struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshToken struct {
	ID        int64
	UserID    int64
	Token     string
	ExpiresAt time.Time
}

type ContextKey string
