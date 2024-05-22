package database

import "errors"

var (
	ErrChirpTooLong  = errors.New("chirp is too long")
	ErrChirpNotFound = errors.New("chirp not found")
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type Database struct {
	Chirps     map[int]Chirp `json:"chirps"`
	Users      map[int]User  `json:"users"`
	NextID     int           `json:"next_id"`
	NextUserID int           `json:"next_user_id"`
}

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}
