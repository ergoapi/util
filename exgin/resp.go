package exgin

import (
	"fmt"

	errors "github.com/ergoapi/util/exerror"
	"github.com/ergoapi/util/ztime"

	"github.com/gin-gonic/gin"
)

type response struct {
	Code      int    `json:"code"`
	Data      any    `json:"data"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
	TraceID   string `json:"traceId"`
}

func newResponse(code int, traceID string, data any, message string) *response {
	return &response{
		Timestamp: ztime.NowUnix(),
		Code:      code,
		Message:   message,
		TraceID:   traceID,
		Data:      data,
	}
}

// getTraceID 从上下文获取追踪ID
func getTraceID(c *gin.Context) string {
	return c.Writer.Header().Get("X-Trace-Id")
}

func SucessResponse(c *gin.Context, data any) {
	traceID := getTraceID(c)
	c.JSON(200, newResponse(200, traceID, data, "请求成功"))
}

func ErrorResponse(c *gin.Context, httpcode int, err error) {
	traceID := getTraceID(c)
	c.JSON(httpcode, newResponse(httpcode, traceID, nil, err.Error()))
}

// ErrorResponse2xx 处理错误响应, 状态码为200
func ErrorResponse2xx(c *gin.Context, code int, err error) {
	traceID := getTraceID(c)
	c.JSON(200, newResponse(code, traceID, nil, err.Error()))
}

// GinsData 处理通用数据响应
func GinsData(c *gin.Context, code int, data any, err error) {
	if err == nil {
		SucessResponse(c, data)
		return
	}
	ErrorResponse(c, code, err)
}

// GinsData2xx 处理通用数据响应, 状态码为200
func GinsData2xx(c *gin.Context, code int, data any, err error) {
	if err == nil {
		SucessResponse(c, data)
		return
	}
	ErrorResponse2xx(c, code, err)
}

// GinsAbort 中止请求并返回错误信息
func GinsAbort(c *gin.Context, httpcode int, msg string) {
	traceID := getTraceID(c)
	c.AbortWithStatusJSON(httpcode, newResponse(httpcode, traceID, nil, msg))
}

// GinsAbort200 中止请求并返回自定义状态码
func GinsAbort200(c *gin.Context, code int, msg string) {
	traceID := getTraceID(c)
	c.AbortWithStatusJSON(200, newResponse(code, traceID, nil, msg))
}

// Bind 绑定JSON请求体
func Bind(c *gin.Context, ptr interface{}) {
	err := c.ShouldBindJSON(ptr)
	if err != nil {
		errors.Bomb("参数不合法: %v", err)
	}
}

// BindWithErr 绑定JSON请求体并返回错误
func BindWithErr(c *gin.Context, ptr interface{}) error {
	err := c.ShouldBindJSON(ptr)
	if err != nil {
		return fmt.Errorf("参数不合法: %v", err)
	}
	return nil
}
