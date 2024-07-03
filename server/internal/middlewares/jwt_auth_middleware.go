package middlewares

import (
	"peer-talk/internal/auth"
	"peer-talk/internal/common"
	"strings"

	"github.com/labstack/echo"
)

func extractTokenFromHeaderOrQueryParam(c echo.Context) (string, error) {
	authHeader := c.Request().Header.Get("Authorization")

	if authHeader == "" {
		authHeader = c.QueryParam("authorization")
	}

	t := strings.Split(authHeader, " ")
	if len(t) != 2 {
		return "", common.NewUnauthorizedError("Unauthorized")
	}

	authToken := t[1]
	authorized, err := auth.IsAuthorized(authToken)

	if err != nil {
		c.Logger().Error(err)
		return "", err
	}

	if !authorized {
		return "", common.NewUnauthorizedError("Unauthorized")
	}

	return authToken, nil
}

func JwtAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authToken, err := extractTokenFromHeaderOrQueryParam(c)

		if err != nil {
			c.Logger().Error(err)
			return err
		}

		authenticatedUser, err := auth.GetAuthUserFromToken(authToken)
		if err != nil {
			return common.NewUnauthorizedError("Unauthorized")

		}

		c.Set(string(common.AuthenticatedUserContextKey), authenticatedUser)

		return next(c)
	}
}
