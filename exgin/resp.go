package exgin

import (
	"fmt"

	errors "github.com/ergoapi/util/exerror"
	"github.com/ergoapi/util/ztime"
	"github.com/gin-gonic/gin"
)

// respDone done
func respDone(code int, data interface{}) gin.H {
	return gin.H{
		"data":      data,
		"message":   "请求成功",
		"timestamp": ztime.NowUnix(),
		"code":      code,
	}
}

// respError error
func respError(code int, data interface{}) gin.H {
	return gin.H{
		"data":      nil,
		"message":   data,
		"timestamp": ztime.NowUnix(),
		"code":      code,
	}
}

func renderMessage(c *gin.Context, v interface{}) {
	if v == nil {
		c.JSON(200, respDone(200, nil))
		return
	}

	switch t := v.(type) {
	case string:
		c.JSON(200, respError(10400, t))
	case error:
		c.JSON(200, respError(10400, t.Error()))
	}
}

func GinsData(c *gin.Context, data interface{}, err error) {
	if err == nil {
		c.JSON(200, respDone(200, data))
		return
	}

	renderMessage(c, err.Error())
}

func GinsCodeData(c *gin.Context, code int, data interface{}, err error) {
	if err == nil {
		c.JSON(200, respDone(code, data))
		return
	}

	renderMessage(c, err.Error())
}

func GinsErrorData(c *gin.Context, code int, data interface{}, err error) {
	c.JSON(200, gin.H{
		"data":      data,
		"message":   fmt.Sprintf("%v", err),
		"timestamp": ztime.NowUnix(),
		"code":      code,
	})
}

func GinsAbort(c *gin.Context, msg string, args ...interface{}) {
	c.AbortWithStatusJSON(200, respError(10400, fmt.Sprintf(msg, args...)))
}

func GinsAbortWithCode(c *gin.Context, respcode int, msg string, args ...interface{}) {
	c.AbortWithStatusJSON(200, respError(respcode, fmt.Sprintf(msg, args...)))
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
