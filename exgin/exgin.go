// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package exgin

import (
	"fmt"
	"time"

	"github.com/ergoapi/util/exnet"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/google/gops/agent"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	promNamespace = "exgin"
	promGinLabels = []string{
		"status_code",
		"url",
		"path",
		"method",
	}
	promGinReqCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: promNamespace,
			Name:      "req_count",
			Help:      "gin server request count",
		}, promGinLabels,
	)
	promGinReqLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: promNamespace,
			Name:      "req_latency",
			Help:      "gin server request latency in seconds",
		}, promGinLabels,
	)
	// 默认慢请求时间 3s
	defaultGinSlowThreshold = time.Second * 3
)

type Config struct {
	Debug          bool
	Gops           bool
	GopsPath       string
	Pprof          bool
	PprofPath      string
	NoCors         bool
	NoTrace        bool
	Metrics        bool
	MetricsPath    string
	TrustedProxies []string
}

func (c *Config) GinSet(r *gin.Engine) {
	gin.DisableConsoleColor()
	if c.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	if len(c.TrustedProxies) > 0 {
		r.SetTrustedProxies(c.TrustedProxies)
	} else {
		r.SetTrustedProxies([]string{"0.0.0.0/0", "::/0"})
	}
	if !c.NoCors {
		r.Use(exCors())
	}

	if !c.NoTrace {
		r.Use(exTraceID())
	}
	if c.Gops {
		if c.GopsPath == "" {
			c.GopsPath = "0.0.0.0:32388"
		}
		go agent.Listen(agent.Options{
			Addr:            c.GopsPath,
			ShutdownCleanup: true})
	}
	if c.Pprof {
		if c.PprofPath == "" {
			c.PprofPath = fmt.Sprintf("/hostdebug/%v/entry", exnet.LocalIPs()[0])
		}
		pprof.Register(r, c.PprofPath)
	}
	if c.Metrics {
		if c.MetricsPath == "" {
			c.MetricsPath = "/metrics"
		}
		r.GET(c.MetricsPath, gin.WrapH(promhttp.Handler()))
	}
}

// Init init gin engine
func Init(c *Config) *gin.Engine {
	r := gin.New()
	r.HandleMethodNotAllowed = true
	c.GinSet(r)
	return r
}

func RealIP(c *gin.Context) string {
	xff := c.Writer.Header().Get("X-Forwarded-For")
	if xff == "" {
		return c.ClientIP()
	}
	return xff
}

func Host(c *gin.Context) string {
	h := c.Request.Host
	if h == "" {
		return c.Request.URL.Host
	}
	return h
}
