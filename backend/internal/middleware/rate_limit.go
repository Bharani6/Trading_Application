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
func RateLimiter(pMaxRequests int, pWindowSeconds int) gin.HandlerFunc {
	lWindow := time.Duration(pWindowSeconds) * time.Second

	return func(pCtx *gin.Context) {
		lIp := pCtx.ClientIP()

		rateLimitMu.Lock()
		lEntry, lExists := rateLimitStore[lIp]

		if !lExists || time.Now().After(lEntry.resetAt) {
			// New window
			rateLimitStore[lIp] = &rateLimitEntry{
				count:   1,
				resetAt: time.Now().Add(lWindow),
			}
			rateLimitMu.Unlock()
			pCtx.Next()
			return
		}

		lEntry.count++
		lCount := lEntry.count
		rateLimitMu.Unlock()

		if lCount > pMaxRequests {
			pCtx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "Too many requests. Please wait and try again.",
			})
			return
		}

		pCtx.Next()
	}
}
