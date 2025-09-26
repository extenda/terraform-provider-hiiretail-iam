package client

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

// LogLevel represents the level of logging
type LogLevel int

const (
	LogLevelNone LogLevel = iota
	LogLevelError
	LogLevelInfo
	LogLevelDebug
)

// Logger handles logging for the client
type Logger struct {
	level LogLevel
	log   *log.Logger
}

// NewLogger creates a new logger with the specified level
func NewLogger(level LogLevel) *Logger {
	return &Logger{
		level: level,
		log:   log.New(os.Stderr, "", log.LstdFlags),
	}
}

// Error logs error messages
func (l *Logger) Error(format string, v ...interface{}) {
	if l.level >= LogLevelError {
		l.log.Printf("[ERROR] "+format, v...)
	}
}

// Info logs informational messages
func (l *Logger) Info(format string, v ...interface{}) {
	if l.level >= LogLevelInfo {
		l.log.Printf("[INFO] "+format, v...)
	}
}

// Debug logs debug messages
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level >= LogLevelDebug {
		l.log.Printf("[DEBUG] "+format, v...)
	}
}

// LogRequest logs HTTP request details
func (l *Logger) LogRequest(req *http.Request) {
	if l.level >= LogLevelDebug {
		dump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			l.Error("Failed to dump request: %v", err)
			return
		}
		l.Debug("HTTP Request:\n%s", string(dump))
	} else if l.level >= LogLevelInfo {
		l.Info("HTTP Request: %s %s", req.Method, req.URL.String())
	}
}

// LogResponse logs HTTP response details
func (l *Logger) LogResponse(resp *http.Response) {
	if l.level >= LogLevelDebug {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			l.Error("Failed to dump response: %v", err)
			return
		}
		l.Debug("HTTP Response:\n%s", string(dump))
	} else if l.level >= LogLevelInfo {
		l.Info("HTTP Response: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
}

// LogRetry logs retry attempts
func (l *Logger) LogRetry(attempt int, err error) {
	if l.level >= LogLevelInfo {
		if err != nil {
			l.Info("Retry attempt %d failed: %v", attempt, err)
		} else {
			l.Info("Retry attempt %d", attempt)
		}
	}
}

// formatRequest formats basic request information for logging
func formatRequest(req *http.Request) string {
	return fmt.Sprintf("%s %s", req.Method, req.URL.String())
}