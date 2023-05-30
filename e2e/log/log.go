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

func (l *Logger) Faillnf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	logrus.Fatalln(fmt.Sprintf("%s: %s", l.prefix, message))
}

func (l *Logger) Printlnf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	logrus.Println(fmt.Sprintf("%s: %s", l.prefix, message))
}

func (l *Logger) Infolnf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	logrus.Infoln(fmt.Sprintf("%s: %s", l.prefix, message))
}
