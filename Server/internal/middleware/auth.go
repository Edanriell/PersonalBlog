package middleware

import (
	"net/http"
	"strings"

	"Server/internal/utils"
)

func JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}

			// Remove "Bearer " prefix
			if strings.HasPrefix(token, "Bearer ") {
				token = token[7:]
			}

			claims, err := utils.ValidateJWT(token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			// Add user ID to context
			c.Set("user_id", claims.UserID)
			return next(c)
		}
	}
}
