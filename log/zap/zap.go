package zap

import (
	"fmt"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/miaocansky/go-tool/util/file"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"time"
)

func NewZap() (logger *zap.Logger) {
	if ok, _ := file.PathExists(zapConfig.Director); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", zapConfig.Director)
		_ = os.Mkdir(zapConfig.Director, os.ModePerm)
	}

	switch zapConfig.Level { // 初始化配置文件的Level
	case "debug":
		zapLevel = zap.DebugLevel
	case "info":
		zapLevel = zap.InfoLevel
	case "warn":
		zapLevel = zap.WarnLevel
	case "error":
		zapLevel = zap.ErrorLevel
	case "dpanic":
		zapLevel = zap.DPanicLevel
	case "panic":
		zapLevel = zap.PanicLevel
	case "fatal":
		zapLevel = zap.FatalLevel
	default:
		zapLevel = zap.InfoLevel
	}

	if zapLevel == zap.DebugLevel || zapLevel == zap.ErrorLevel {
		logger = zap.New(getEncoderCore(), zap.AddStacktrace(zapLevel))
	} else {
		logger = zap.New(getEncoderCore())
	}
	if zapConfig.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}

//
//  getEncoderConfig
//  @Description:获取配置
//  @return config
//
func getEncoderConfig() (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:    "message",
		LevelKey:      "level",
		TimeKey:       "time",
		NameKey:       "logger",
		CallerKey:     "caller",
		StacktraceKey: zapConfig.StacktraceKey,
		LineEnding:    zapcore.DefaultLineEnding,
		//EncodeLevel:   zapcore.CapitalColorLevelEncoder,
		EncodeLevel: customEncodeLevel,

		EncodeTime: customTimeEncoder,
		//EncodeTime:     cEncodeTime,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		//EncodeDuration: secondsDurationEncoder,
		EncodeCaller: zapcore.FullCallerEncoder,
	}
	//switch {
	//case zapUtil.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
	//	config.EncodeLevel = zapcore.LowercaseLevelEncoder
	//case zapUtil.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
	//	config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	//case zapUtil.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
	//	config.EncodeLevel = zapcore.CapitalLevelEncoder
	//case zapUtil.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
	//	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	//default:
	//	config.EncodeLevel = zapcore.LowercaseLevelEncoder
	//}
	return config
}

// getEncoder 获取zapcore.Encoder
func getEncoder() zapcore.Encoder {
	if zapConfig.Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

// getEncoderCore 获取Encoder的zapcore.Core
func getEncoderCore() (core zapcore.Core) {
	writer, err := GetWriteSyncer() // 使用file-rotatelogs进行日志分割
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return
	}
	return zapcore.NewCore(getEncoder(), writer, zapLevel)
}

// 日志写入方式
func GetWriteSyncer() (zapcore.WriteSyncer, error) {
	fileWriter, err := rotatelogs.New(
		path.Join(zapConfig.Director, "%Y-%m-%d.log"),
		rotatelogs.WithLinkName(zapConfig.LinkName),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if zapConfig.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), err
	}
	return zapcore.AddSync(fileWriter), err
}

// 日志写入方式2
func GetWriteSyncer2() (zapcore.WriteSyncer, error) {
	// 日志文件 与 日志切割 配置
	lumberJackLogger := &lumberjack.Logger{
		Filename:   path.Join(zapConfig.Director, "test.log"), // 日志文件路径
		MaxSize:    1,                                         // 单个日志文件最大多少 mb
		MaxBackups: 10,                                        // 日志备份数量
		MaxAge:     7,                                         // 日志最长保留时间
		Compress:   false,                                     // 是否压缩日志
	}

	if zapConfig.LogInConsole {
		// 日志同时输出到控制台和日志文件中
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberJackLogger), zapcore.AddSync(os.Stdout)), nil
	} else {
		// 日志只输出到日志文件
		return zapcore.AddSync(lumberJackLogger), nil
	}
}

// 自定义时间格式
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(zapConfig.Prefix + "[" + "2006/01/02 - 15:04:05.000" + "]"))
}

// customEncodeLevel 自定义日志级别显示
func customEncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}
