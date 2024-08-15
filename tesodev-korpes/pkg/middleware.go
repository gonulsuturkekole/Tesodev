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
		// List of paths to skip
		skipPaths := []string{"/login", "/verify", "/customer"}
		// Get the request path
		reqPath := c.Path()
		// Check if the request path should be skipped
		for _, path := range skipPaths {
			if strings.HasPrefix(reqPath, path) {
				return next(c) // Skip the middleware
			}
		}
		tokenString := c.Request().Header.Get("Authentication")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "No Authentication header provided"})
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Check if the token is a valid JWT
		if err := authentication.VerifyJWT(tokenString); err != nil {
			return err
		}
		// Call the verify endpoint with the token
		verifyUrl := "http://localhost:8001/verify"
		req, err := http.NewRequest("GET", verifyUrl, nil)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create verification request"})
		}
		req.Header.Set("Authentication", tokenString)

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Verification request failed"})
		}
		defer res.Body.Close()
		//if res.StatusCode != http.StatusOK {
		//	return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token verification failed"})
		//}

		return next(c)
	}
}
