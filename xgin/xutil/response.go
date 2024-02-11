package xutil

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Data      any    `json:"data"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestId string `json:"request_id"`
}

func newResponse(c *gin.Context, code int, data any, msg string) *Response {
	response := &Response{
		Code:      code,
		Message:   msg,
		Data:      data,
		RequestId: GetTraceId(c),
	}
	c.Set(ResponseBodyMsgKey, response.Message)
	c.Set(ResponseBodyCodeKey, response.Code)
	return response
}

func DefaultResponse(c *gin.Context) *Response {
	return newResponse(c, 0, nil, "")
}

func NewResponse(c *gin.Context, code int, data any, msg string) *Response {
	return newResponse(c, code, data, msg)
}

func SystemErrorResponse(c *gin.Context) *Response {
	response := newResponse(c, SystemError, nil, "system error")
	return response
}

func BodyValidateErrorResponse(c *gin.Context, msg string) *Response {
	response := newResponse(c, RequestBodyValidateError, nil, msg)
	return response
}

func FailedResponse(c *gin.Context, msg string, codeList ...int) *Response {
	code := OtherError
	if len(codeList) != 0 {
		code = codeList[0]
	}
	response := newResponse(c, code, nil, msg)
	return response
}

func SuccessResponse(c *gin.Context, data any, msgList ...string) *Response {
	msg := ""
	if len(msgList) != 0 {
		msg = msgList[0]
	}
	response := newResponse(c, NoError, data, msg)
	return response
}
