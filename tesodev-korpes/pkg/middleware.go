package pkg

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	_ "github.com/labstack/echo/v4"
	"net/http"
	"strings"
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
		skipPaths := []string{"/login"}
		reqPath := c.Path()

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

		// Verify token and handle logic
		isVerified, err := VerifyTokenWithAPI(tokenString)
		if err != nil || !isVerified {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid or expired token"})
		}

		return next(c)
	}
}

// VerifyTokenWithAPI sends a request to the verify endpoint
func VerifyTokenWithAPI(token string) (bool, error) {
	// Create request payload
	payload := map[string]string{"token": token}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return false, err
	}

	// Send POST request to verify endpoint
	resp, err := http.Post("http://localhost:8001/verify", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	// Handle response body if needed (e.g., read JSON response)
	return true, nil
}
