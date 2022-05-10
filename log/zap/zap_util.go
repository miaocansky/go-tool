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
func (zapUtil *ZapUtil) Error(msg string, loggerData dto.LoggerData) {
	zapUtil.Zap.Error(msg, zap.Any(loggerData.Key, loggerData.Data))
}
func (zapUtil *ZapUtil) Info(msg string, loggerData dto.LoggerData) {
	zapUtil.Zap.Info(msg, zap.Any(loggerData.Key, loggerData.Data))
}
func (zapUtil *ZapUtil) Debug(msg string, loggerData dto.LoggerData) {
	zapUtil.Zap.Debug(msg, zap.Any(loggerData.Key, loggerData.Data))

}
func (zapUtil *ZapUtil) Warn(msg string, loggerData dto.LoggerData) {
	zapUtil.Zap.Warn(msg, zap.Any(loggerData.Key, loggerData.Data))
}
