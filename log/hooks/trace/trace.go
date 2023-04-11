package trace

import "github.com/sirupsen/logrus"

type TraceIdHook struct {
	TraceID string
}

func NewTraceIdHook(traceID string) logrus.Hook {
	hook := TraceIdHook{
		TraceID: traceID,
	}
	return &hook
}

func (hook *TraceIdHook) Fire(entry *logrus.Entry) error {
	entry.Data["traceID"] = hook.TraceID
	entry.Data["Tag"] = "gin"
	return nil
}

func (hook *TraceIdHook) Levels() []logrus.Level {
	return logrus.AllLevels
}