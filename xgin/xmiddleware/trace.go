package xmiddleware

import (
	"github.com/blackRice-Tu/golib/xgin/xutil"
	"github.com/gin-gonic/gin"
)

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// inject trace id
		xutil.InjectTraceId(c)
		c.Next()
		return
	}
}
