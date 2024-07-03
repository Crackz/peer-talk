package auth

type AuthUser struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,notblank,min=1,max=255"`
	Username string `json:"username" validate:"required,notblank,min=1,max=255"`
	Password string `json:"password" validate:"required,notblank,min=4,max=255"`
	Email    string `json:"email" validate:"omitempty,email"`
}

type RegisterResponse struct {
	AccessToken string    `json:"accessToken"`
	User        *AuthUser `json:"user"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=1"`
	Password string `json:"password" validate:"required,min=4"`
}

type LoginResponse struct {
	AccessToken string    `json:"accessToken"`
	User        *AuthUser `json:"user"`
}
