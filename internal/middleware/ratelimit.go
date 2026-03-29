package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	lastSeen time.Time
	count    int
}

var (
	visitors        = &sync.Map{}
	mu              sync.Mutex
	lastCleanup     time.Time
	cleanupInterval = 5 * time.Minute
)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		cleanupIfNeeded()

		ip := c.ClientIP()

		now := time.Now()
		value, loaded := visitors.LoadOrStore(ip, &visitor{lastSeen: now, count: 1})
		v := value.(*visitor)

		if !loaded {
			c.Next()
			return
		}

		mu.Lock()
		if now.Sub(v.lastSeen) > time.Minute {
			v.count = 1
			v.lastSeen = now
			mu.Unlock()
			c.Next()
			return
		}

		if v.count > 50 {
			mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please try again later.",
			})
			return
		}

		v.count++
		v.lastSeen = now
		mu.Unlock()

		c.Next()
	}
}

func cleanupIfNeeded() {
	mu.Lock()
	defer mu.Unlock()

	if time.Since(lastCleanup) < cleanupInterval {
		return
	}
	lastCleanup = time.Now()

	threshold := time.Now().Add(-cleanupInterval)

	visitors.Range(func(key, value any) bool {
		v := value.(*visitor)
		if v.lastSeen.Before(threshold) {
			visitors.Delete(key)
		}
		return true
	})
}
