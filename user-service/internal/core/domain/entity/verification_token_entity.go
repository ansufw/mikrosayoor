package entity

import "time" // Asumsi Anda memerlukan import time jika ada field time.Time

type VerificationTokenEntity struct {
	ID        int64
	UserID    int64
	Token     string
	TokenType string
	ExpiresAt time.Time
	User      UserEntity
}
