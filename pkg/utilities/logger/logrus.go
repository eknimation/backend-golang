package logger

import (
	"github.com/sirupsen/logrus"
	"go.elastic.co/ecslogrus"
)

func GetLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&ecslogrus.Formatter{})
	return logger
}
