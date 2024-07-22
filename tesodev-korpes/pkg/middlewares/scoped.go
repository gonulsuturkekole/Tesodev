package middlewares

import (
	"github.com/labstack/echo/v4"
	"strconv"
	"time"
)

var scopedServices = make(map[string]*ScopedService)

type (
	ScopedService struct {
		Id string
	}
)

func generateUniqueID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

// ScopedServiceMiddleware is the middleware function.
func ScopedServiceMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestId := generateUniqueID()
		scopedService := &ScopedService{
			Id: "",
		}
		scopedServices[requestId] = scopedService

		if err := next(c); err != nil {
			c.Error(err)
		}

		delete(scopedServices, requestId)
		return nil
	}
}
