package logger

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	L *logrus.Logger
}

func (l *LogrusLogger) Log(c echo.Context, logType string, msg ...interface{}) {
	fields := logrus.Fields{
		"Request Method": c.Request().Method,
		"Remote Address": c.Request().RemoteAddr,
		"Message":        msg,
	}

	switch logType {
	case "error":
		l.L.WithFields(fields).Error(c.Request().URL.Path)
	case "info":
		l.L.WithFields(fields).Info(c.Request().URL.Path)
	case "warning":
		l.L.WithFields(fields).Warning(c.Request().URL.Path)
	}
}

func NewLogrusLogger() *LogrusLogger {
	return &LogrusLogger{
		L: logrus.New(),
	}
}
