// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package main

import (
	"errors"
	"net/http"
	"os"

	utilerrors "github.com/ergoapi/util/exerror"
	"github.com/ergoapi/util/exgin"
	"github.com/ergoapi/util/exhttp"
	"github.com/ergoapi/util/exid"
	filehook "github.com/ergoapi/util/log/hooks/file"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.AddHook(filehook.NewRotateFileHook(filehook.RotateFileConfig{
		Filename:   "/tmp/ergoapi.log",
		MaxSize:    10,
		MaxBackups: 1,
		MaxAge:     1,
		Level:      logrus.DebugLevel,
		Formatter:  &logrus.JSONFormatter{},
	}))
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
}

func main() {
	g := exgin.Init(&exgin.Config{
		Debug: true,
	})
	g.Use(exgin.ExLog(), exgin.ExRecovery())
	g.GET("/", func(ctx *gin.Context) {
		nextid := exid.GenSnowflakeID()
		exgin.SucessResponse(ctx, map[string]any{
			"ip": exgin.RealIP(ctx),
			"snowid": map[string]any{
				"id":    nextid,
				"parse": exid.ParseID(nextid),
			},
		})
	})
	g.POST("/admin", func(ctx *gin.Context) {
		exgin.GinsData(ctx, 200, nil, nil)
	})
	g.POST("/panic", func(ctx *gin.Context) {
		panic("panic")
	})
	g.POST("/panic2", func(ctx *gin.Context) {
		utilerrors.Bomb("参数不合法")
	})
	g.NoMethod(func(c *gin.Context) {
		exgin.GinsData2xx(c, 400, nil, errors.New("not support method"))
	})
	g.NoRoute(func(c *gin.Context) {
		exgin.GinsData(c, 400, nil, errors.New("not found route"))
	})
	addr := "0.0.0.0:65001"
	srv := &http.Server{
		Addr:    addr,
		Handler: g,
	}
	go func() {
		exhttp.SetupGracefulStop(srv)
	}()
	logrus.Infof("http listen to %v, pid is %v", addr, os.Getpid())
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Errorf("Failed to start http server, error: %s", err)
	}
}
