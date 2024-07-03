package auth

import (
	"fmt"
	"peer-talk/config"
	"peer-talk/internal/user"
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtCustomClaims struct {
	jwt.StandardClaims
	Name string `json:"name,omitempty"`
}

type JwtCustomRefreshClaims struct {
	jwt.StandardClaims
	Name string `json:"name,omitempty"`
}

type AuthenticatedUser struct {
	Id   string
	Name string
}

func CreateAccessToken(user *user.User) (accessToken string, err error) {
	expirationInSeconds := config.Env.JwtAccessTokenExpirationInSeconds
	expiredAt := time.Now().Add(time.Duration(expirationInSeconds) * time.Second).Unix()
	claims := &JwtCustomClaims{
		Name: user.Name,
		StandardClaims: jwt.StandardClaims{
			Subject:   user.ID,
			ExpiresAt: expiredAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(config.Env.JwtAccessTokenSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func CreateRefreshToken(user *user.User) (refreshToken string, err error) {
	expirationInSeconds := config.Env.JwtRefreshTokenExpirationInSeconds
	expiredAt := time.Now().Add(time.Duration(expirationInSeconds) * time.Second).Unix()
	claimsRefresh := &JwtCustomRefreshClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   user.ID,
			ExpiresAt: expiredAt,
		},
		Name: user.Name,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)

	signedToken, err := token.SignedString([]byte(config.Env.JwtRefreshTokenSecret))
	if err != nil {
		return "", err
	}

	return signedToken, err
}

func IsAuthorized(requestToken string) (bool, error) {
	_, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Env.JwtAccessTokenSecret), nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetAuthUserFromToken(requestToken string) (*AuthenticatedUser, error) {
	token, err := jwt.ParseWithClaims(requestToken, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Env.JwtAccessTokenSecret), nil
	})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	claims, ok := token.Claims.(*JwtCustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid Token")
	}

	return &AuthenticatedUser{
		Id:   claims.Subject,
		Name: claims.Name,
	}, nil
}
