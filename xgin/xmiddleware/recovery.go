package xmiddleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"

	"github.com/blackRice-Tu/golib/xgin/xutil"
)

// RecoveryMiddleware ...
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				errorInfo := fmt.Sprintf("%s\n%s", errorToString(r), string(debug.Stack()))
				c.Set(xutil.ErrorInfoKey, errorInfo)
				resp := xutil.SystemErrorResponse(c)
				c.AbortWithStatusJSON(http.StatusOK, resp)
			}
		}()
		c.Next()
	}
}

// recover错误，转string
func errorToString(r interface{}) string {
	switch v := r.(type) {
	case error:
		return v.Error()
	default:
		return r.(string)
	}
}
