package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/ergoapi/util/exgin"
	"github.com/ergoapi/util/exhttp"
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
}

func main() {
	g := exgin.Init(&exgin.Config{
		Debug: true,
	})
	g.Use(exgin.ExLog(), exgin.ExRecovery())
	g.NoRoute(func(c *gin.Context) {
		exgin.GinsErrorData(c, 404, nil, errors.New("not found route"))
	})
	g.Any("/", func(ctx *gin.Context) {
		exgin.GinsData(ctx, nil, nil)
	})
	g.NoMethod(func(c *gin.Context) {
		exgin.GinsErrorData(c, 404, nil, errors.New("not support method"))
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
