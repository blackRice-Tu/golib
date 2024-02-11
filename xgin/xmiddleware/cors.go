package xmiddleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setCorsHeader(c *gin.Context, origin string) {
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Origin", origin)
	c.Header("Access-Control-Allow-Headers", "content-type,start-client-header")
}

// CorsMiddleware ...
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if len(origin) == 0 {
			return
		}

		host := c.Request.Host
		if origin == "http://"+host || origin == "https://"+host {
			return
		}

		setCorsHeader(c, origin)
		if c.Request.Method == "OPTIONS" {
			defer c.AbortWithStatus(http.StatusNoContent)
		}
	}
}
