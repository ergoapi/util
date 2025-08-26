// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

// Package trace provides tracing hooks.
package trace

import "github.com/sirupsen/logrus"

type IDHook struct {
	TraceID    string
	TraceAgent string
}

func NewTraceIDHook(traceID string, TraceAgent ...string) logrus.Hook {
	hook := IDHook{
		TraceID: traceID,
	}
	if len(TraceAgent) > 0 {
		hook.TraceAgent = TraceAgent[0]
	} else {
		hook.TraceAgent = "ergoapi-sdk"
	}
	return &hook
}

func (hook *IDHook) Fire(entry *logrus.Entry) error {
	entry.Data["traceID"] = hook.TraceID
	entry.Data["Tag"] = hook.TraceAgent
	return nil
}

func (hook *IDHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
