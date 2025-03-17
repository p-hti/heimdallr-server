package model

import "time"

type User struct {
	ID       int64  `db:"id"`
	Email    string `db:"email"`
	Name     string `db:"name"`
	PassHash []byte `db:"pass_hash"`
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
	CreatedAt time.Time
}

type ContextKey string
