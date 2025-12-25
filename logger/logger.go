package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func New() *Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})
	log.SetLevel(logrus.InfoLevel)

	if os.Getenv("DEBUG") == "true" {
		log.SetLevel(logrus.DebugLevel)
	}

	return &Logger{Logger: log}
}

func (l *Logger) WithContext(fields map[string]interface{}) *logrus.Entry {
	return l.WithFields(logrus.Fields(fields))
}
