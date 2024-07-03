package user

import (
	db "peer-talk/db/sqlc"
)

type CreatableUser struct {
	Name           string
	Username       string
	Email          string
	HashedPassword string
}

type User struct {
	ID       string
	Name     string
	Username string
	Email    string
	Password string
}

type Store interface {
	CreateUser(createUser CreatableUser) (db.User, error)
	GetUserByUsername(username string) (db.User, error)
}
