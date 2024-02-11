package xmiddleware

import (
	"bytes"
	"github.com/blackRice-Tu/golib"
	"github.com/blackRice-Tu/golib/utils/xcommon"
	"github.com/blackRice-Tu/golib/xgin/xutil"
	"github.com/blackRice-Tu/golib/xlogger"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

var (
	loggerCollector *zap.Logger
)

const (
	CollectResultSucceeded   = "Succeeded"
	CollectResultFailed      = "Failed"
	CollectResultClientError = "ClientError"
	CollectResultServerError = "ServerError"
	CollectResultOther       = "Other"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

type collectRequestBody struct {
	// request info
	ClientIp string `json:"client_ip"`
	Uri      string `json:"uri"`
	Scheme   string `json:"scheme"`
	Host     string `json:"host"`
	Header   string `json:"header"`
	Method   string `json:"method"`
	Query    string `json:"query"`
	Body     string `json:"body"`

	// server info
	Env        string `json:"env"`
	ServerIp   string `json:"server_ip"`
	ProgramKey string `json:"program_key"`
	TraceId    string `json:"trace_id"`
	ErrorInfo  string `json:"error_info"`

	// response info
	StatusCode       int    `json:"status_code"`
	ResponseBody     string `json:"response_body"`
	ResponseBodyCode int    `json:"response_body_code"`
	ResponseBodyMsg  string `json:"response_body_msg"`

	// analyze
	Result string `json:"result"`

	// custom info
	Custom map[string]any `json:"custom"`
}

type CollectRequestOpt struct {
	EscapeFunc     func(ctx *gin.Context) bool
	CustomInfoFunc func(ctx *gin.Context) map[string]any
	LoggerConfig   *xlogger.Config
}

func CollectRequestMiddleware(opt *CollectRequestOpt) gin.HandlerFunc {
	if opt == nil {
		opt = &CollectRequestOpt{}
	}

	loggerCollector = xlogger.NewLogger(opt.LoggerConfig)

	return func(c *gin.Context) {
		if opt.EscapeFunc != nil && opt.EscapeFunc(c) {
			c.Next()
			return
		}

		collectorBody := collectRequestBody{}

		// request info
		collectorBody.ClientIp = c.ClientIP()
		collectorBody.Uri = c.Request.URL.Path
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		collectorBody.Scheme = scheme
		collectorBody.Host = c.Request.Host
		collectorBody.Header = xcommon.JsonMarshal(c.Request.Header)
		collectorBody.Method = c.Request.Method
		collectorBody.Query = xcommon.JsonMarshal(c.Request.URL.Query())
		collectorBody.Body, _ = xutil.GetRequestBody(c)

		// get response content
		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w
		c.Next()

		// server info
		collectorBody.Env = golib.GetEnv()
		collectorBody.ServerIp = golib.GetServerIp()
		collectorBody.ProgramKey = golib.GetProgramKey()
		collectorBody.TraceId = xutil.GetTraceId(c)
		errorInfo, ok := c.Get(xutil.ErrorInfoKey)
		if ok {
			collectorBody.ErrorInfo = errorInfo.(string)
		}

		// response info
		collectorBody.StatusCode = c.Writer.Status()
		collectorBody.ResponseBody = string(w.body.Bytes())
		if responseCode, ok := c.Get(xutil.ResponseBodyCodeKey); ok {
			collectorBody.ResponseBodyCode = responseCode.(int)
		}
		if responseMsg, ok := c.Get(xutil.ResponseBodyMsgKey); ok {
			collectorBody.ResponseBodyMsg = responseMsg.(string)
		}

		// custom info
		if opt.CustomInfoFunc != nil {
			collectorBody.Custom = opt.CustomInfoFunc(c)
		}

		// result
		result := CollectResultOther
		if collectorBody.StatusCode >= 500 {
			result = CollectResultServerError
		} else if collectorBody.StatusCode >= 400 && collectorBody.StatusCode < 500 {
			result = CollectResultClientError
		} else if collectorBody.StatusCode >= 200 && collectorBody.StatusCode < 300 {
			if collectorBody.ResponseBodyCode == 0 {
				result = CollectResultSucceeded
			} else {
				result = CollectResultFailed
			}
		}
		collectorBody.Result = result

		// collect to logger
		if loggerCollector != nil {
			loggerCollector.Info(xcommon.JsonMarshal(collectorBody))
		}
	}
}
