package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	prefix string
}

func NewLogger(prefix string) *Logger {
	return &Logger{prefix: prefix}
}

func (l *Logger) Failln(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	logrus.Fatalln(fmt.Sprintf("%s: %s", l.prefix, message))
}

func (l *Logger) Println(args ...interface{}) {
	logrus.Println(fmt.Sprintf("%s: %s", l.prefix, fmt.Sprint(args...)))
}

func (l *Logger) Infoln(args ...interface{}) {
	logrus.Infoln(fmt.Sprintf("%s: %s", l.prefix, fmt.Sprint(args...)))
}
