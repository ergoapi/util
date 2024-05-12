package exgin

import (
	"fmt"

	errors "github.com/ergoapi/util/exerror"
	"github.com/ergoapi/util/ztime"
	"github.com/gin-gonic/gin"
)

// customRespDone done
func customRespDone(code int, traceId string, data, message interface{}) gin.H {
	msg := message
	switch t := message.(type) {
	case error:
		msg = t.Error()
	}
	return gin.H{
		"data":      data,
		"message":   msg,
		"timestamp": ztime.NowUnix(),
		"code":      code,
		"traceId":   traceId,
		"success":   code == 200,
	}
}

// respDone done
func respDone(code int, traceId string, data interface{}) gin.H {
	return customRespDone(code, traceId, data, "请求成功")
}

// respError error
func respError(code int, traceId string, message interface{}) gin.H {
	return customRespDone(code, traceId, nil, message)
}

func renderMessage(c *gin.Context, traceId string, v interface{}) {
	if v == nil {
		c.JSON(200, respDone(200, traceId, nil))
		return
	}

	switch t := v.(type) {
	case string:
		c.JSON(200, respError(10400, traceId, t))
	case error:
		c.JSON(200, respError(10400, traceId, t.Error()))
	}
}

func GinsData(c *gin.Context, data interface{}, err error) {
	traceId := c.Writer.Header().Get("X-Trace-Id")
	if err == nil {
		c.JSON(200, respDone(200, traceId, data))
		return
	}
	renderMessage(c, traceId, err.Error())
}

func GinsCodeData(c *gin.Context, code int, data interface{}, err error) {
	traceId := c.Writer.Header().Get("X-Trace-Id")
	if err == nil {
		c.JSON(200, respDone(code, traceId, data))
		return
	}
	renderMessage(c, traceId, err.Error())
}

func GinsErrorData(c *gin.Context, code int, data interface{}, err error) {
	traceId := c.Writer.Header().Get("X-Trace-Id")
	c.JSON(200, customRespDone(code, traceId, data, err))
}

func GinsAbort(c *gin.Context, msg string, args ...interface{}) {
	traceId := c.Writer.Header().Get("X-Trace-Id")
	c.AbortWithStatusJSON(200, respError(10400, traceId, fmt.Sprintf(msg, args...)))
}

func GinsAbortWithCode(c *gin.Context, respcode int, msg string, args ...interface{}) {
	traceId := c.Writer.Header().Get("X-Trace-Id")
	c.AbortWithStatusJSON(200, respError(respcode, traceId, fmt.Sprintf(msg, args...)))
}

func GinsCustomResp(c *gin.Context, obj interface{}) {
	c.JSON(200, obj)
}

func GinsCustomCodeResp(c *gin.Context, code int, obj interface{}) {
	c.JSON(code, obj)
}

func Bind(c *gin.Context, ptr interface{}) {
	err := c.ShouldBindJSON(ptr)
	if err != nil {
		errors.Bomb("参数不合法: %v", err)
	}
}

// BindWithErr bind with error
func BindWithErr(c *gin.Context, ptr interface{}) error {
	err := c.ShouldBindJSON(ptr)
	if err != nil {
		return fmt.Errorf("参数不合法: %v", err)
	}
	return nil
}

// APICustomRespBody swag api resp body
type APICustomRespBody struct {
	Code      int         `json:"code"`
	Data      interface{} `json:"data"`
	Message   string      `json:"message"`
	Timestamp int         `json:"timestamp"`
}
