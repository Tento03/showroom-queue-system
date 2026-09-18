package middleware

import (
	"backend-queue/config"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RateLimitUpload() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Key berdasarkan IP
		ip := c.ClientIP()
		key := "ratelimit:upload:" + ip

		ctx := context.Background()

		// Cek berapa kali sudah request
		count, err := config.RDB.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}

		// Set expire hanya di request pertama
		if count == 1 {
			config.RDB.Expire(ctx, key, 10*time.Second)
		}

		// Max 5 request per 10 detik per IP
		if count > 5 {
			ttl, _ := config.RDB.TTL(ctx, key).Result()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "too many requests",
				"retry_after": ttl.Seconds(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
