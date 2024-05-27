package database

import "errors"
import "time"

var (
	ErrChirpTooLong  = errors.New("chirp is too long")
	ErrChirpNotFound = errors.New("chirp not found")
	ErrUserNotFound  = errors.New("user not found")
	ErrUserExists    = errors.New("user already exists")
)

type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	AuthorID int    `json:"author_id"`
}

type Database struct {
	Chirps        map[int]Chirp           `json:"chirps"`
	Users         map[int]User            `json:"users"`
	NextID        int                     `json:"next_id"`
	NextUserID    int                     `json:"next_user_id"`
	RefreshTokens map[string]RefreshToken `json:"refresh_tokens"`
}

type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type RefreshToken struct {
	Token     string    `json:"token"`
	UserID    int       `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}
