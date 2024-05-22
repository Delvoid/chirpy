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
	Chirps map[int]Chirp `json:"chirps"`
	NextID int           `json:"next_id"`
}
