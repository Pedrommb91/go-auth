package models

import (
	"time"
)

type Credentials struct {
	ID        int32     `name:"id"`
	Salt      string    `name:"salt"`
	PassHash  string    `name:"passhash"`
	CreatedAt time.Time `name:"created_at"`
	UpdatedAt time.Time `name:"updated_at"`
}

type Users struct {
	ID          int32       `name:"id"`
	Username    string      `name:"username"`
	Email       string      `name:"email"`
	Credentials Credentials `name:"credentials_id" reference:"credentials"`
	CreatedAt   time.Time   `name:"created_at"`
	UpdatedAt   time.Time   `name:"updated_at"`
}

type UserReaderInterface interface {
}

type UserWriterInterface interface {
	AddUser(user Users) (int64, error)
}

type UserRepositoryInterface interface {
	UserReaderInterface
	UserWriterInterface
}
