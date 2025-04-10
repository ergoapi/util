package exgin

import (
	"strconv"

	errors "github.com/ergoapi/util/exerror"

	"github.com/gin-gonic/gin"
)

// GinsQueryStr query string
func GinsQueryStr(c *gin.Context, key string, defaultval ...string) string {
	val := c.Query(key)
	if val != "" {
		return val
	}
	if len(defaultval) == 0 {
		errors.Bomb("query param[%s] is necessary", key)
	}
	return defaultval[0]
}

// GinsQueryStrNull query string
func GinsQueryStrNull(c *gin.Context, key string) string {
	val := c.Query(key)
	if val != "" {
		return val
	}
	return ""
}

// GinsQueryInt query int
func GinsQueryInt(c *gin.Context, key string, defaultVal ...int) int {
	strv := c.Query(key)
	if strv != "" {
		intv, err := strconv.Atoi(strv)
		if err != nil {
			errors.Bomb("cannot convert [%s] to int", strv)
		}
		return intv
	}

	if len(defaultVal) == 0 {
		errors.Bomb("query param[%s] is necessary", key)
	}

	return defaultVal[0]
}

// GinsQueryInt64 querry int64
func GinsQueryInt64(c *gin.Context, key string, defaultVal ...int64) int64 {
	strv := c.Query(key)
	if strv != "" {
		intv, err := strconv.ParseInt(strv, 10, 64)
		if err != nil {
			errors.Bomb("cannot convert [%s] to int64", strv)
		}
		return intv
	}

	if len(defaultVal) == 0 {
		errors.Bomb("query param[%s] is necessary", key)
	}

	return defaultVal[0]
}

// GinsQueryBool query bool
func GinsQueryBool(c *gin.Context, key string, defaultVal ...bool) bool {
	strv := c.Query(key)
	if strv != "" {
		intv, err := strconv.Atoi(strv)
		if err != nil {
			return false
		}
		return intv == 1
	}

	if len(defaultVal) == 0 {
		return false
	}

	return defaultVal[0]
}

// GinsParamStr param str
func GinsParamStr(c *gin.Context, field string) string {
	val := c.Param(field)
	if val == "" {
		errors.Bomb("url param[%s] is null", field)
	}
	return val
}

// GinsParamInt64 param str
func GinsParamInt64(c *gin.Context, field string) int64 {
	strval := GinsParamStr(c, field)
	intval, err := strconv.ParseInt(strval, 10, 64)
	if err != nil {
		errors.Bomb("cannot convert %s to int64", strval)
	}
	return intval
}

// GinsParamInt param str
func GinsParamInt(c *gin.Context, field string) int {
	return int(GinsParamInt64(c, field))
}

// GinsOffset offset
func GinsOffset(c *gin.Context, limit int) int {
	if limit <= 0 {
		limit = 10
	}
	page := GinsQueryInt(c, "page", 1)
	return (page - 1) * limit
}

// GinsHeader header key
func GinsHeader(c *gin.Context, headerkey string) string {
	return c.GetHeader(headerkey)
}
