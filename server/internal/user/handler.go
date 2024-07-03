package user

import (
	"database/sql"
	db "peer-talk/db/sqlc"
	"peer-talk/internal/common"
)

type Handler struct {
	store Store
}

func NewHandler(store Store) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) mapToUser(user db.User) *User {
	mappedUser := &User{
		ID:       user.ID.String(),
		Name:     user.Name,
		Username: user.Username,
		Password: user.Password,
	}

	if user.Email.Valid {
		mappedUser.Email = user.Email.String
	}
	return mappedUser

}

func (h *Handler) CreateUser(creatableUser CreatableUser) (*User, error) {
	existedUser, err := h.store.GetUserByUsername(creatableUser.Username)
	if err != nil && err != sql.ErrNoRows {
		return nil, common.NewInternalServerError(err.Error())

	}

	isUsernameExisted := existedUser.Username != ""
	if isUsernameExisted {
		return nil, common.NewConflictError(common.ApiErrorDetails{
			Message: "username already exists",
			Param:   "userExists",
		})
	}

	createdUser, err := h.store.CreateUser(creatableUser)
	if err != nil {
		return nil, err
	}

	return h.mapToUser(createdUser), nil
}

func (h *Handler) GetUserByUsername(username string) (*User, error) {
	existedUser, err := h.store.GetUserByUsername(username)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if existedUser.Username == "" {
		return nil, nil
	}

	return h.mapToUser(existedUser), nil
}
