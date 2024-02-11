package xrequests

import (
	"context"
	"net/http"
	"reflect"
	"time"

	"github.com/blackRice-Tu/golib"
	"github.com/blackRice-Tu/golib/utils/xcommon"
	"github.com/blackRice-Tu/golib/utils/xcontext"
	logger "github.com/blackRice-Tu/golib/xlogger/default"

	"github.com/levigross/grequests"
	"github.com/pkg/errors"
)

type (
	RequestOptions = grequests.RequestOptions
	Response       = grequests.Response

	SendFunc func(ctx context.Context, url string, params any, responseStruct any, opts ...*RequestOptions) (trace *Trace, response *Response, e error)
)

const (
	DefaultTimeoutSecond   = 15
	DefaultKeepAliveSecond = 5
)

type Trace struct {
	Url          string            `json:"url"`
	Method       string            `json:"method"`
	Headers      map[string]string `json:"headers"`
	Query        any               `json:"query"`
	Body         any               `json:"body"`
	StatusCode   int               `json:"status_code"`
	ResponseBody string            `json:"response_body"`
}

func buildRequestOptions(ctx context.Context, opts []*RequestOptions) *RequestOptions {
	opt := &RequestOptions{}
	for _, o := range opts {
		if o != nil {
			opt = o
			break
		}
	}
	if opt.DialTimeout == 0 {
		opt.DialTimeout = time.Duration(DefaultTimeoutSecond) * time.Second
	}
	if opt.RequestTimeout == 0 {
		opt.RequestTimeout = time.Duration(DefaultTimeoutSecond) * time.Second
	}
	if opt.DialKeepAlive == 0 {
		opt.DialKeepAlive = time.Duration(DefaultKeepAliveSecond) * time.Second
	}
	// inject traceId
	traceIdKey := golib.GetTraceIdKey()
	if opt.Headers == nil {
		opt.Headers = make(map[string]string)
	}
	if _, ok := opt.Headers[traceIdKey]; !ok {
		opt.Headers[traceIdKey] = xcontext.GetOrNewTraceId(ctx)
	}
	return opt
}

func newTrace(ctx context.Context, url string, method string, opt *RequestOptions) *Trace {
	trace := &Trace{
		Url:     url,
		Method:  method,
		Headers: opt.Headers,
	}
	if opt.Data != nil {
		trace.Body = opt.Data
	} else if opt.JSON != nil {
		trace.Body = opt.JSON
	}
	trace.Query = opt.Params
	return trace
}

func wrapResponse(ctx context.Context, trace *Trace, response *Response, err error) (*Trace, *Response, error) {
	responseBody := string(response.Bytes())
	trace.StatusCode = response.StatusCode
	trace.ResponseBody = responseBody

	if xcontext.IsDebug(ctx) {
		logger.Debugf(ctx, "[Trace] %s", xcommon.JsonMarshal(trace))
	}

	if err != nil {
		return trace, response, err
	}
	if !response.Ok {
		err = errors.Errorf("StatusCode: %d; Body: %s", response.StatusCode, responseBody)
		return trace, response, err
	}
	return trace, response, err
}

func NativeGet(ctx context.Context, url string, opts ...*RequestOptions) (trace *Trace, response *Response, e error) {
	opt := buildRequestOptions(ctx, opts)
	trace = newTrace(ctx, url, http.MethodGet, opt)
	defer func() {
		if e != nil {
			logger.Errorf(ctx, "[Err] %s; [Trace] %s", e.Error(), xcommon.JsonMarshal(trace))
		}
	}()
	response, e = grequests.Get(url, buildRequestOptions(ctx, opts))
	trace, response, e = wrapResponse(ctx, trace, response, e)
	return
}

func NativePost(ctx context.Context, url string, opts ...*RequestOptions) (trace *Trace, response *Response, e error) {
	opt := buildRequestOptions(ctx, opts)
	trace = newTrace(ctx, url, http.MethodPost, opt)
	defer func() {
		if e != nil {
			logger.Errorf(ctx, "[Err] %s; [Trace] %s", e.Error(), xcommon.JsonMarshal(trace))
		}
	}()
	response, e = grequests.Post(url, opt)
	trace, response, e = wrapResponse(ctx, trace, response, e)
	return
}

func Get(ctx context.Context, url string, params any, responseStruct any, opts ...*RequestOptions) (trace *Trace, response *Response, e error) {
	queryMap, err := buildQuery(params)
	if err != nil {
		e = errors.WithMessagef(err, "buildQuery")
		return
	}
	opt := buildRequestOptions(ctx, opts)
	opt.Params = queryMap
	trace, response, e = NativeGet(ctx, url, opt)
	if e != nil {
		return
	}
	err = jsonUnmarshalResponse(response, responseStruct)
	if err != nil {
		e = errors.Wrapf(err, "response.JSON")
	}
	return
}

func Post(ctx context.Context, url string, params any, responseStruct any, opts ...*RequestOptions) (trace *Trace, response *Response, e error) {
	opt := buildRequestOptions(ctx, opts)
	opt.JSON = params
	trace, response, e = NativePost(ctx, url, opt)
	if e != nil {
		return
	}
	err := jsonUnmarshalResponse(response, responseStruct)
	if err != nil {
		e = errors.Wrapf(err, "response.JSON")
	}
	return
}

func jsonUnmarshalResponse(response *Response, v any) (e error) {
	if v == nil {
		return
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer {
		e = errors.Errorf("non-pointer field")
		return
	}
	e = response.JSON(v)
	return
}
