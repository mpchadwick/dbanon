package dbanon

import (
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
}

func SetLogger(l *logrus.Logger) {
	logger = l
}

func GetLogger() *logrus.Logger {
	return logger
}