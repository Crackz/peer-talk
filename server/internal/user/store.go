package user

import (
	"context"
	"database/sql"
	db "peer-talk/db/sqlc"
)

type store struct {
	sqlStore *db.SQLStore
	ctx      context.Context
}

func NewStore(sqlStore *db.SQLStore) *store {
	return &store{sqlStore: sqlStore, ctx: context.Background()}
}

func (s *store) CreateUser(creatableUser CreatableUser) (db.User, error) {
	var email sql.NullString
	if creatableUser.Email != "" {
		email.String = creatableUser.Email
		email.Valid = true
	}

	return s.sqlStore.CreateUser(s.ctx, db.CreateUserParams{
		Name:     creatableUser.Name,
		Username: creatableUser.Username,
		Email:    email,
		Password: creatableUser.HashedPassword,
	})
}

func (s *store) GetUserByUsername(username string) (db.User, error) {
	return s.sqlStore.GetUserByUsername(s.ctx, username)

}
