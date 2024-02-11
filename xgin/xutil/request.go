package xutil

import (
	"bytes"
	"io"

	"github.com/blackRice-Tu/golib"
	"github.com/blackRice-Tu/golib/utils/xcontext"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

// GetRequestBody ...
func GetRequestBody(c *gin.Context) (string, error) {
	body, err := c.GetRawData()
	if err != nil {
		return "", errors.Wrapf(err, "GetRawData")
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	return string(body), nil
}

func InjectTraceId(c *gin.Context) {
	traceId := c.Request.Header.Get(golib.GetTraceIdKey())
	if traceId == "" {
		traceId = xcontext.GenerateTraceId()
	}
	c.Set(golib.GetTraceIdKey(), traceId)
	return
}

func GetTraceId(c *gin.Context) (traceId string) {
	traceId = xcontext.GetOrNewTraceId(c)
	return
}

func BindRequestBody(c *gin.Context, body any) (e error) {
	err := c.ShouldBindBodyWith(body, binding.JSON)
	if err != nil {
		e = errors.Wrapf(err, "json.Unmarshal err")
		return
	}
	validate := validator.New()
	err = validate.Struct(body)
	if err != nil {
		e = errors.Wrapf(err, "validate err")
		return
	}
	return
}
