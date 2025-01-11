package logging

import (
	"context"

	"github.com/sirupsen/logrus"
)

type ContextKey int

const (
	// RequestIDKey is a custom context key for the request ID.
	RequestIDKey ContextKey = 0
)

var Log = logrus.New()

// InitLogger initializes the logger with desired configurations
func InitLogger() {
	// Set logrus to use the JSON formatter
	Log.SetFormatter(&logrus.JSONFormatter{})
}

// logWithFields logs a message with custom fields and the specified log level
func logWithFields(level logrus.Level, fields logrus.Fields, message string, args ...interface{}) {
	entry := Log.WithFields(fields)
	switch level {
	case logrus.InfoLevel:
		if len(args) > 0 {
			entry.Infof(message, args...)
		} else {
			entry.Info(message)
		}
	case logrus.WarnLevel:
		if len(args) > 0 {
			entry.Warnf(message, args...)
		} else {
			entry.Warn(message)
		}
	case logrus.ErrorLevel:
		if len(args) > 0 {
			entry.Errorf(message, args...)
		} else {
			entry.Error(message)
		}
	default:
		entry.Info(message)
	}
}

// LogCustomField logs a message with custom fields
func LogCustomField(level logrus.Level, fields logrus.Fields, message string, args ...interface{}) {
	logWithFields(level, fields, message, args...)
}

// LogInfo logs an informational message with request ID from context
func LogInfo(ctx context.Context, message string, args ...interface{}) {
	logID, _ := ctx.Value(RequestIDKey).(string)
	LogCustomField(logrus.InfoLevel, logrus.Fields{"request_id": logID}, message, args...)
}

// LogWarning logs a warning message with request ID from context
func LogWarning(ctx context.Context, message string, args ...interface{}) {
	logID, _ := ctx.Value(RequestIDKey).(string)
	LogCustomField(logrus.WarnLevel, logrus.Fields{"request_id": logID}, message, args...)
}

// LogError logs an error message with request ID from context
func LogError(ctx context.Context, message string, args ...interface{}) {
	logID, _ := ctx.Value(RequestIDKey).(string)
	LogCustomField(logrus.ErrorLevel, logrus.Fields{"request_id": logID}, message, args...)
}
