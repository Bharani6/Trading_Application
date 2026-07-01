package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimitEntry struct {
	count   int
	resetAt time.Time
}

var (
	rateLimitStore = make(map[string]*rateLimitEntry)
	rateLimitMu    sync.Mutex
)

// RateLimiter returns a middleware that limits requests per IP.
// maxRequests: max number of requests allowed within windowSeconds.
func RateLimiter(maxRequests int, windowSeconds int) gin.HandlerFunc {
	window := time.Duration(windowSeconds) * time.Second

	return func(c *gin.Context) {
		ip := c.ClientIP()

		rateLimitMu.Lock()
		entry, exists := rateLimitStore[ip]

		if !exists || time.Now().After(entry.resetAt) {
			// New window
			rateLimitStore[ip] = &rateLimitEntry{
				count:   1,
				resetAt: time.Now().Add(window),
			}
			rateLimitMu.Unlock()
			c.Next()
			return
		}

		entry.count++
		count := entry.count
		rateLimitMu.Unlock()

		if count > maxRequests {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "Too many requests. Please wait and try again.",
			})
			return
		}

		c.Next()
	}
}
