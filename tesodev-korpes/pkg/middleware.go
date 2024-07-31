package pkg

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	_ "github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"tesodev-korpes/CustomerService/authentication"
)

func CorrelationIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Checking X-Correlation-Id header
		correlationID := c.Request().Header.Get("X-Correlation-Id")
		// if not exist , generate new UUID
		if correlationID == "" {
			correlationID = uuid.New().String()
			// Add new UUID to header
			c.Request().Header.Set("X-Correlation-Id", correlationID)
		}
		// Add response header too
		c.Response().Header().Set("X-Correlation-Id", correlationID)
		// continue
		return next(c)
	}
}

func Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authentication")
		if tokenString == "" {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "No Authentication header provided"})
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		_, err := authentication.VerifyJWT(tokenString)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		//username is chosen inside claims
		//c.Set("_id", claims.Id)
		return next(c)
	}
}
