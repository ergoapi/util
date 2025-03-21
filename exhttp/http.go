//  Copyright (c) 2021. The EFF Team Authors.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  See the License for the specific language governing permissions and
//  limitations under the License.

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

// SetupGracefulStop grace stop
func SetupGracefulStop(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logrus.Infof("receive signal %s", sig)
	ShutDown(srv)
}

// ShutDown http shutdown
func ShutDown(srv *http.Server) {
	logrus.Info("service stopping...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Errorf("service stop failed: %v", err)
	}
	select {
	case <-ctx.Done():
		logrus.Warn("service stop timeout")
	default:
		logrus.Info("service stopped")
	}
	logrus.Info("server exited.")
}
