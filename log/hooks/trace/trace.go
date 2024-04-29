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
