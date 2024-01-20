package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	log *logrus.Logger
}

func New() *Logger {
	l := logrus.New()
	l.Out = os.Stdout
	l.Formatter = &logrus.JSONFormatter{}
	l.Level = logrus.DebugLevel

	return &Logger{
		log: l,
	}
}

func (l *Logger) LogWithFields(
	lvl logrus.Level,
	msg string,
	attributes map[string]string,
) {
	fields := logrus.Fields{}

	// Put specific attributes
	for key, val := range attributes {
		fields[key] = val
	}

	// Put common attributes
	fields["service.name"] = "client"

	switch lvl {
	case logrus.ErrorLevel:
		l.log.WithFields(fields).Error(msg)
	case logrus.InfoLevel:
		l.log.WithFields(fields).Info(msg)
	default:
		l.log.WithFields(fields).Debug(msg)
	}
}
