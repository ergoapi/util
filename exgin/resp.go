package exgin

import (
	"fmt"

	errors "github.com/ergoapi/util/exerror"
	"github.com/ergoapi/util/ztime"
	"github.com/gin-gonic/gin"
)

// customRespDone done
func customRespDone(code int, message, tid, data interface{}) gin.H {
	return gin.H{
		"data":      data,
		"message":   message,
		"timestamp": ztime.NowUnix(),
		"code":      code,
		"traceId":   tid,
		"success":   code == 200,
	}
}

// respDone done
func respDone(code int, tid, data interface{}) gin.H {
	return customRespDone(code, "请求成功", tid, data)
}

// respError error
func respError(code int, message, tid interface{}) gin.H {
	return customRespDone(code, message, tid, nil)
}

func renderMessage(c *gin.Context, v interface{}) {
	tid := c.Writer.Header().Get("X-Trace-Id")
	if v == nil {
		c.JSON(200, respDone(200, tid, nil))
		return
	}

	switch t := v.(type) {
	case string:
		c.JSON(200, respError(10400, tid, t))
	case error:
		c.JSON(200, respError(10400, tid, t.Error()))
	}
}

func GinsData(c *gin.Context, data interface{}, err error) {
	tid := c.Writer.Header().Get("X-Trace-Id")
	if err == nil {
		c.JSON(200, respDone(200, tid, data))
		return
	}

	renderMessage(c, err.Error())
}

func GinsCodeData(c *gin.Context, code int, data interface{}, err error) {
	tid := c.Writer.Header().Get("X-Trace-Id")
	if err == nil {
		c.JSON(200, respDone(code, tid, data))
		return
	}

	renderMessage(c, err.Error())
}

func GinsErrorData(c *gin.Context, code int, data interface{}, err error) {
	tid := c.Writer.Header().Get("X-Trace-Id")
	c.JSON(200, customRespDone(code, fmt.Sprintf("%v", err), tid, data))
}

func GinsAbort(c *gin.Context, msg string, args ...interface{}) {
	tid := c.Writer.Header().Get("X-Trace-Id")
	c.AbortWithStatusJSON(200, respError(10400, fmt.Sprintf(msg, args...), tid))
}

func GinsAbortWithCode(c *gin.Context, respcode int, msg string, args ...interface{}) {
	tid := c.Writer.Header().Get("X-Trace-Id")
	c.AbortWithStatusJSON(200, respError(respcode, fmt.Sprintf(msg, args...), tid))
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

// APICustomRespBody swag api resp body
type APICustomRespBody struct {
	Code      int         `json:"code"`
	Data      interface{} `json:"data"`
	Message   string      `json:"message"`
	Timestamp int         `json:"timestamp"`
}
