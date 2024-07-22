package middlewares

import (
	"github.com/labstack/echo/v4"
	"sync"
	"time"
)

// Singleton instance
var stats *Stats

type (
	Stats struct {
		Uptime       time.Time `json:"uptime"`
		RequestCount uint64    `json:"requestCount"`
		mutex        sync.RWMutex
	}
)

func NewStats() *Stats {
	if stats == nil {
		stats = &Stats{
			Uptime: time.Now(),
		}
	}
	return stats
}

// Process is the middleware function.
func (s *Stats) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			c.Error(err)
		}
		s.mutex.Lock()
		defer s.mutex.Unlock()
		s.RequestCount++
		return nil
	}
}
