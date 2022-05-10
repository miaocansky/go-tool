package log

import "github.com/miaocansky/go-tool/log/dto"

type Logger interface {
	Error(msg string, loggerData dto.LoggerData)
	Info(msg string, loggerData dto.LoggerData)
	Debug(msg string, loggerData dto.LoggerData)
	Warn(msg string, loggerData dto.LoggerData)
}
