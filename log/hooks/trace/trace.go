package trace

import "github.com/sirupsen/logrus"

type IDHook struct {
	TraceID string
}

func NewTraceIDHook(traceID string) logrus.Hook {
	hook := IDHook{
		TraceID: traceID,
	}
	return &hook
}

func (hook *IDHook) Fire(entry *logrus.Entry) error {
	entry.Data["traceID"] = hook.TraceID
	entry.Data["Tag"] = "gin"
	return nil
}

func (hook *IDHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
