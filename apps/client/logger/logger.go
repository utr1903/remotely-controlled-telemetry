package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	log *logrus.Logger
}

func New() *Logger {
	l := logrus.New()
	l.Out = os.Stdout

	logsDir := "./logs"
	logFile := "log"

	// Create the directory if it doesn't exist
	err := os.MkdirAll(logsDir, 0700)
	if err != nil {
		fmt.Println("Directory creation failed.")
		panic(err)
	}

	// Add permissions
	err = os.Chmod(logsDir, 0700)
	if err != nil {
		fmt.Println("Error changing directory permissions:", err)
	}

	// Create the file in the "logs" directory
	filePath := filepath.Join(logsDir, logFile)
	_, err = os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		panic(err)
	}
	// defer file.Close()

	// Open the file and pass it to logrus
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	l.Out = file

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
