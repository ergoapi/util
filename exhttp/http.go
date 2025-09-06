// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package exhttp

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

// SetupGracefulStop registers signal handlers and gracefully shuts down the server when a signal is received.
// Note: This call blocks; run it in a separate goroutine, e.g.:
//
//	go SetupGracefulStop(srv)
func SetupGracefulStop(srv *http.Server) {
	SetupGracefulStopWithTimeout(srv, 5*time.Second)
}

// SetupGracefulStopWithTimeout is like SetupGracefulStop but allows customizing shutdown timeout and signals.
// If no signals are provided, it defaults to SIGINT and SIGTERM.
func SetupGracefulStopWithTimeout(srv *http.Server, timeout time.Duration, sigs ...os.Signal) {
	quit := make(chan os.Signal, 1)
	if len(sigs) == 0 {
		sigs = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}
	signal.Notify(quit, sigs...)
	defer signal.Stop(quit)

	sig := <-quit
	logrus.Infof("receive signal %s", sig)
	ShutDownWithTimeout(srv, timeout)
}

// ShutDown gracefully shuts down the server with a default 5s timeout.
func ShutDown(srv *http.Server) {
	ShutDownWithTimeout(srv, 5*time.Second)
}

// ShutDownWithTimeout gracefully shuts down the server with the provided timeout.
// If graceful shutdown fails, it will attempt to force close the server.
func ShutDownWithTimeout(srv *http.Server, timeout time.Duration) {
	logrus.Info("service stopping...")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Errorf("service stop failed: %v", err)
		_ = srv.Close()
	}

	if ctx.Err() == context.DeadlineExceeded {
		logrus.Warn("service stop timeout")
	} else {
		logrus.Info("service stopped")
	}
	logrus.Info("server exited.")
}
