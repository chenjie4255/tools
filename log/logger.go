package log

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var defaultLogger *Logger

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	defaultLogger = NewLogger("")
}

// SetDefaultLogger 设置默认的logger
func SetDefaultLogger(topic string) {
	defaultLogger = NewLogger(topic)
}

// Default 获取默认logger
func Default() *Logger {
	return defaultLogger
}

// Logger logger
type Logger struct {
	*logrus.Entry
	enableSentry bool
}

// Fields 结构化输出的Map type
type Fields logrus.Fields

// NewLogger 新增一个logger
func NewLogger(topic string) *Logger {
	entry := logrus.WithField("topic", topic)
	return &Logger{entry, false}
}

// NewLoggerWithSentry 新增一个logger,带sentry报错功能
func NewLoggerWithSentry(topic string) *Logger {
	entry := logrus.WithField("topic", topic)
	return &Logger{entry, true}
}

// SetLogLevel 设置日志等级
func SetLogLevel(level string) error {
	l, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.Errorf("Parse level error, %s, %s", level, err)
		return err
	}
	logrus.SetLevel(l)
	return nil
}

// WithFields 批量添加Key value数据
func (l *Logger) WithFields(fields Fields) *Logger {
	entry := l.Entry.WithFields(logrus.Fields(fields))
	return &Logger{entry, l.enableSentry}
}

// AddFile 添加调用文件信息
func (l *Logger) AddFile() *Logger {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return l
	}

	loc := fmt.Sprintf("%s:%d", file, line)
	splitIndex := strings.LastIndex(loc, "moreless.space/mobile/")
	if splitIndex == -1 {
		return l
	}

	entry := l.WithField("file", loc[splitIndex+17:])
	return &Logger{entry, l.enableSentry}
}

func (l *Logger) buildSentryTags() map[string]string {
	ret := make(map[string]string)
	for k, v := range l.Data {
		ret[k] = fmt.Sprintf("%v", v)
	}
	return ret
}

func (l *Logger) Warn(args ...interface{}) {
	l.Entry.Warn(args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Entry.Warnf(format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.Entry.Error(args...)
}

// InfoError 添加调用文件信息
func (l *Logger) InfoError(err error, args ...interface{}) {
	if err != nil {
		l.WithField("error", err).Error(args...)
		return
	}

	l.Info(args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Entry.Errorf(format, args...)
}
