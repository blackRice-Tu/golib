package xutil

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func JsonRender(c *gin.Context, response *Response) {
	c.JSON(http.StatusOK, response)
}
