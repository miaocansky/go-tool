package log

import (
	"github.com/miaocansky/go-tool/log/zap"
	"strconv"
	"testing"
)

func TestRegister(t *testing.T) {
	//zap.NewZapUtil()
	config := zap.ZapConfig{
		Director:      "log",
		Level:         "info",
		ShowLine:      true,
		StacktraceKey: "S",
		Format:        "console",
		Prefix:        "[GIN-WEB]",
		LinkName:      "latest_log",
		LogInConsole:  false,
	}

	Register("zap", zap.NewZapUtil(config))
	i := 1
	//for i := 0; i < 5000; i++ {
	DefaultDriver().Info("onemsg")
	DefaultDriver().Info("msg", Ang("test=>", "msg:"+strconv.Itoa(i)))
	DefaultDriver().Error("msg", Ang("test=>", "msg:"+strconv.Itoa(i)))
	DefaultDriver().Debug("msg", Ang("test=>", "msg:"+strconv.Itoa(i)))
	DefaultDriver().Warn("msg", Ang("test=>", "msg:"+strconv.Itoa(i)))
	//}

}
