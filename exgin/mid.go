// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package exgin

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/ergoapi/util/environ"
	errors "github.com/ergoapi/util/exerror"
	"github.com/ergoapi/util/exid"
	ltrace "github.com/ergoapi/util/log/hooks/trace"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/sirupsen/logrus"
)

var uni = ut.New(en.New(), zh.New())

// exCors ex cors middleware
func exCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, UPDATE, HEAD, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, X-Auth-Token, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Access-Control-Request-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Max-Age", "3600")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func exTraceID() gin.HandlerFunc {
	return func(g *gin.Context) {
		traceID := g.GetHeader("X-Trace-Id")
		if traceID == "" {
			traceID = exid.GenUUID()
			g.Request.Header.Set("X-Trace-Id", traceID)
		}
		g.Header("X-Trace-Id", traceID)
		g.Set("ex-trace-id", traceID)
		logrus.AddHook(ltrace.NewTraceIDHook(traceID, "exgin"))
		g.Next()
	}
}

// ExLog log middleware
func ExLog(skip ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		host := Host(c)
		path := c.Request.URL.Path
		method := c.Request.Method
		ua := c.Request.UserAgent()
		query := c.Request.URL.RawQuery
		c.Next()
		for _, s := range skip {
			if strings.HasPrefix(path, s) {
				return
			}
		}
		end := time.Now()
		latency := end.Sub(start)
		if len(query) == 0 {
			query = " - "
		}
		if latency > defaultGinSlowThreshold {
			logrus.Warnf("[msg] api %v query %v", path, latency)
		}
		statuscode := c.Writer.Status()
		bodysize := c.Writer.Size()
		clientIP := c.ClientIP()
		remoteIP := c.RemoteIP()
		xffIP := c.Writer.Header().Get("X-Forwarded-For")
		readIP := c.Writer.Header().Get("X-Real-Ip")
		referer := c.Request.Referer()
		if len(c.Errors) > 0 || c.Writer.Status() >= 500 {
			logrus.WithFields(logrus.Fields{
				"statuscode": statuscode,
				"bodysize":   bodysize,
				"client_ip":  clientIP,
				"remote_ip":  remoteIP,
				"xff_ip":     xffIP,
				"real_ip":    readIP,
				"method":     method,
				"host":       host,
				"path":       path,
				"latency":    latency,
				"ua":         ua,
				"referer":    referer,
			}).Warnf("query: %v  <= err: %v", query, c.Errors.String())
		} else {
			logrus.WithFields(logrus.Fields{
				"statuscode": statuscode,
				"bodysize":   bodysize,
				"client_ip":  clientIP,
				"remote_ip":  remoteIP,
				"xff_ip":     xffIP,
				"real_ip":    readIP,
				"method":     method,
				"host":       host,
				"path":       path,
				"latency":    latency,
				"ua":         ua,
				"referer":    referer,
			}).Infof("query: %v", query)
		}
		// update prom
		labels := []string{fmt.Sprint(statuscode), host, path, method}
		promGinReqCount.WithLabelValues(labels...).Inc()
		promGinReqLatency.WithLabelValues(labels...).Observe(latency.Seconds())
	}
}

// ExRecovery recovery
func ExRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if res, ok := err.(errors.ErgoError); ok {
					code := 400
					if strings.Contains(res.Message, "unauth") {
						code = 401
					}
					GinsAbort(c, code, res.Message)
					return
				}
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logrus.Errorf("Recovery from brokenPipe ---> path: %v, err: %v, request: %v",
						c.Request.URL.Path, err, string(httpRequest))
					GinsAbort(c, 500, "请求broken")
				} else {
					logrus.Errorf("Recovery from panic ---> err: %v, request: %v, stack: %v",
						err, string(httpRequest), string(debug.Stack()))
					GinsAbort(c, 500, "请求panic")
				}
				return
			}
		}()
		c.Next()
	}
}

// ExHackHeader hack header
func ExHackHeader() gin.HandlerFunc {
	return func(g *gin.Context) {
		g.Writer.Header().Add("ex-glb", environ.GetEnv("POD_NAME", "tbh-9526"))
		g.Writer.Header().Add("ex-loc", environ.GetEnv("POD_LOC", "cn"))
		g.Writer.Header().Add("Server", "cloudflare")
		g.Next()
	}
}

// Translations .
func Translations() gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := c.GetHeader("locale")
		trans, _ := uni.GetTranslator(locale)
		v, ok := binding.Validator.Engine().(*validator.Validate)
		if ok {
			switch locale {
			case "zh":
				_ = zh_translations.RegisterDefaultTranslations(v, trans)
			case "en":
				_ = en_translations.RegisterDefaultTranslations(v, trans)
			default:
				_ = zh_translations.RegisterDefaultTranslations(v, trans)
			}
			c.Set("trans", trans)
		}
		c.Next()
	}
}

// NoCache is a middleware function that appends headers
// to prevent the client from caching the HTTP response.
func NoCache(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
	c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	c.Next()
}
