package glog

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ergoapi/util/exctx"
	"github.com/ergoapi/util/file"
	"github.com/ergoapi/util/ztime"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

var DefaultGLogger = GLogger{
	LogLevel:      logger.Info,
	SlowThreshold: 200 * time.Millisecond,
}

type GLogger struct {
	LogLevel      logger.LogLevel
	SlowThreshold time.Duration
}

func (mgl *GLogger) LogMode(logLevel logger.LogLevel) logger.Interface {
	mgl.LogLevel = logLevel
	return mgl
}

func (mgl *GLogger) Info(ctx context.Context, message string, values ...interface{}) {
	trace := exctx.GetTraceContext(ctx)
	msg := fmt.Sprintf("message=%+v||values=%+v", message, fmt.Sprint(values...))
	msg = strings.Trim(fmt.Sprintf("%q", msg), "\"")
	logrus.WithFields(logrus.Fields{
		"traceID":     trace.TraceID,
		"SpanID":      trace.SpanID,
		"childSpanID": trace.CSpanID,
		"Tag":         "gorm",
	}).Info(msg)
}

func (mgl *GLogger) Warn(ctx context.Context, message string, values ...interface{}) {
	trace := exctx.GetTraceContext(ctx)
	msg := fmt.Sprintf("message=%+v||values=%+v", message, fmt.Sprint(values...))
	msg = strings.Trim(fmt.Sprintf("%q", msg), "\"")
	logrus.WithFields(logrus.Fields{
		"traceID":     trace.TraceID,
		"SpanID":      trace.SpanID,
		"childSpanID": trace.CSpanID,
		"Tag":         "gorm",
	}).Warn(msg)
}

func (mgl *GLogger) Error(ctx context.Context, message string, values ...interface{}) {
	trace := exctx.GetTraceContext(ctx)
	msg := fmt.Sprintf("message=%+v||values=%+v", message, fmt.Sprint(values...))
	msg = strings.Trim(fmt.Sprintf("%q", msg), "\"")
	logrus.WithFields(logrus.Fields{
		"traceID":     trace.TraceID,
		"SpanID":      trace.SpanID,
		"childSpanID": trace.CSpanID,
		"Tag":         "gorm",
	}).Error(msg)
}

func (mgl *GLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	trace := exctx.GetTraceContext(ctx)
	if mgl.LogLevel > 0 {
		elapsed := time.Since(begin)
		currentTime := begin.Format(time.RFC3339)
		switch {
		case err != nil && mgl.LogLevel >= logger.Error:
			sql, rows := fc()
			if rows == -1 || rows == 0 || err == gorm.ErrRecordNotFound {
				go file.WritefileWithLine(fmt.Sprintf("%v.dbnotfound.txt", ztime.GetToday()), sql)
				logrus.WithFields(logrus.Fields{
					"traceID":         trace.TraceID,
					"SpanID":          trace.SpanID,
					"childSpanID":     trace.CSpanID,
					"Tag":             "gorm",
					"FileWithLineNum": utils.FileWithLineNum(),
					"current_time":    currentTime,
					"proc_time":       float64(elapsed.Milliseconds()),
					"rows":            "-",
					"sql":             sql,
				}).Warn(err)
			} else {
				logrus.WithFields(logrus.Fields{
					"traceID":         trace.TraceID,
					"SpanID":          trace.SpanID,
					"childSpanID":     trace.CSpanID,
					"Tag":             "gorm",
					"FileWithLineNum": utils.FileWithLineNum(),
					"current_time":    currentTime,
					"proc_time":       float64(elapsed.Milliseconds()),
					"rows":            rows,
					"sql":             sql,
				}).Error(err)
			}
		case mgl.SlowThreshold != 0 && elapsed > mgl.SlowThreshold && mgl.LogLevel >= logger.Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", mgl.SlowThreshold)
			go file.WritefileWithLine(fmt.Sprintf("%v.slowsql.txt", ztime.GetToday()), sql+" "+slowLog)
			if rows == -1 {
				logrus.WithFields(logrus.Fields{
					"traceID":         trace.TraceID,
					"SpanID":          trace.SpanID,
					"childSpanID":     trace.CSpanID,
					"Tag":             "gorm",
					"FileWithLineNum": utils.FileWithLineNum(),
					"current_time":    currentTime,
					"proc_time":       float64(elapsed.Milliseconds()),
					"rows":            "-",
					"sql":             sql,
					"slowlog":         slowLog,
				}).Warn(err)
			} else {
				logrus.WithFields(logrus.Fields{
					"traceID":         trace.TraceID,
					"SpanID":          trace.SpanID,
					"childSpanID":     trace.CSpanID,
					"Tag":             "gorm",
					"FileWithLineNum": utils.FileWithLineNum(),
					"current_time":    currentTime,
					"proc_time":       float64(elapsed.Milliseconds()),
					"rows":            rows,
					"sql":             sql,
					"slowlog":         slowLog,
				}).Warn(err)
			}
		case mgl.LogLevel >= logger.Info:
			sql, rows := fc()
			if rows == -1 {
				logrus.WithFields(logrus.Fields{
					"traceID":         trace.TraceID,
					"SpanID":          trace.SpanID,
					"childSpanID":     trace.CSpanID,
					"Tag":             "gorm",
					"FileWithLineNum": utils.FileWithLineNum(),
					"current_time":    currentTime,
					"proc_time":       float64(elapsed.Milliseconds()),
					"rows":            "-",
					"sql":             sql,
				}).Info(err)
			} else {
				logrus.WithFields(logrus.Fields{
					"traceID":         trace.TraceID,
					"SpanID":          trace.SpanID,
					"childSpanID":     trace.CSpanID,
					"Tag":             "gorm",
					"FileWithLineNum": utils.FileWithLineNum(),
					"current_time":    currentTime,
					"proc_time":       float64(elapsed.Milliseconds()),
					"rows":            rows,
					"sql":             sql,
				}).Info(err)
			}
		}
	}
}
