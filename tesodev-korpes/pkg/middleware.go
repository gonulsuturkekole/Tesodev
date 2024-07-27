package pkg

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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

/*func VerifyToken(c echo.Context) string {

	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Authorization header is required"})
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	return tokenString
}
*/
