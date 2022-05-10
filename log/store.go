package log

import "github.com/miaocansky/go-tool/log/dto"

type Logger interface {
	Error(msg string, loggerDataLists ...dto.LoggerData)
	Info(msg string, loggerDataLists ...dto.LoggerData)
	Debug(msg string, loggerDataLists ...dto.LoggerData)
	Warn(msg string, loggerDataLists ...dto.LoggerData)
}
