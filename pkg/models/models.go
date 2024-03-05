package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord           = errors.New("models: no snippet")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate emial")
)

type News struct {
	ID       string    `bson:"_id,omitempty"`
	Title    string    `bson:"title"`
	Body     string    `bson:"body"`
	ImageURL string    `bson:"imageurl"`
	Created  time.Time `bson:"created"`
}

type User struct {
	ID             string `bson:"_id,omitempty"`
	Name           string `bson:"name"`
	Email          string `bson:"email"`
	HashedPassword []byte `bson:"hashed_password"`
	Role           string `bson:"role"`
}

type Comment struct {
	ID     string `bson:"_id,omitempty"`
	UserID string `bson:"user_id,omitempty"`
	NewsID string `bson:"news_id,omitempty"`
	Text   string `bson:"text"`
}
