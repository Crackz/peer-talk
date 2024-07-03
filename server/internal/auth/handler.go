package auth

import (
	"net/http"
	"peer-talk/internal/common"
	"peer-talk/internal/user"

	"github.com/labstack/echo"
)

const (
	WrongUsernameOrPasswordErrorMsg = "Username or password is wrong"
)

type Handler struct {
	userHandler *user.Handler
}

func NewHandler(userHandler *user.Handler) *Handler {
	return &Handler{userHandler: userHandler}
}

func (h *Handler) RegisterRoutes(router *echo.Group) {
	router.POST("/register", h.handleRegister)
	router.POST("/login", h.handleLogin)
}

func (h *Handler) mapToAuthUser(user *user.User) *AuthUser {
	return &AuthUser{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
	}
}

func (h *Handler) handleRegister(c echo.Context) error {
	var registerRequest RegisterRequest

	if err := common.BindAndValidate(c, &registerRequest); err != nil {
		return err
	}

	hashedPassword, err := HashPassword(registerRequest.Password)
	if err != nil {
		return common.NewInternalServerError(err.Error())
	}

	user, err := h.userHandler.CreateUser(user.CreatableUser{
		Name:           registerRequest.Name,
		Username:       registerRequest.Username,
		Email:          registerRequest.Email,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		return err
	}

	token, err := CreateAccessToken(user)
	if err != nil {
		return common.NewInternalServerError(err.Error())
	}

	c.JSON(http.StatusCreated, &RegisterResponse{
		AccessToken: token,
		User:        h.mapToAuthUser(user),
	})

	return nil
}

func (h *Handler) handleLogin(c echo.Context) error {
	var loginRequest LoginRequest

	if err := common.BindAndValidate(c, &loginRequest); err != nil {
		return err
	}

	user, err := h.userHandler.GetUserByUsername(loginRequest.Username)
	if err != nil {
		return err
	}

	if user == nil {
		return common.NewUnauthorizedError(WrongUsernameOrPasswordErrorMsg)
	}

	if err := CheckPassword(loginRequest.Password, user.Password); err != nil {
		return common.NewUnauthorizedError(WrongUsernameOrPasswordErrorMsg)
	}

	token, err := CreateAccessToken(user)
	if err != nil {
		return common.NewInternalServerError(err.Error())
	}

	c.JSON(http.StatusCreated, &LoginResponse{
		AccessToken: token,
		User:        h.mapToAuthUser(user),
	})

	return nil
}
