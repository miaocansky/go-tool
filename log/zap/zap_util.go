package zap

import (
	"github.com/miaocansky/go-tool/log/dto"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapLevel zapcore.Level
var zapConfig ZapConfig

type ZapUtil struct {
	Zap *zap.Logger
}

func NewZapUtil(config ZapConfig) *ZapUtil {
	zapConfig = ZapConfig{
		Director:      config.Director,
		Level:         config.Level,
		ShowLine:      config.ShowLine,
		StacktraceKey: config.StacktraceKey,
		Format:        config.Format,
		Prefix:        config.Prefix,
		LinkName:      config.LinkName,
		LogInConsole:  config.LogInConsole,
	}

	util := &ZapUtil{}
	util.Zap = NewZap()
	return util
}
func (zapUtil *ZapUtil) Error(msg string, loggerDataLists ...dto.LoggerData) {

	zapUtil.Zap.Error(msg, ExchangeDataToFields(loggerDataLists)...)
}
func (zapUtil *ZapUtil) Info(msg string, loggerDataLists ...dto.LoggerData) {
	zapUtil.Zap.Info(msg, ExchangeDataToFields(loggerDataLists)...)
}
func (zapUtil *ZapUtil) Debug(msg string, loggerDataLists ...dto.LoggerData) {
	zapUtil.Zap.Debug(msg, ExchangeDataToFields(loggerDataLists)...)

}
func (zapUtil *ZapUtil) Warn(msg string, loggerDataLists ...dto.LoggerData) {
	zapUtil.Zap.Warn(msg, ExchangeDataToFields(loggerDataLists)...)
}

func ExchangeDataToFields(loggerDataLists []dto.LoggerData) []zap.Field {
	fields := make([]zap.Field, 0, 8)
	for _, data := range loggerDataLists {
		fields = append(fields, zap.Any(data.Key, data.Value))
	}
	return fields
}
