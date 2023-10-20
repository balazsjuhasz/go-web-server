package database

import "time"

type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	AuthorID int    `json:"author_id"`
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	ID        string    `json:"id"`
	RevokedAt time.Time `json:"revoked_at"`
}
