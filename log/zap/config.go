package zap

type ZapConfig struct {
	Director      string
	Level         string
	ShowLine      bool
	StacktraceKey string
	EncodeLevel   string
	Format        string
	Prefix        string
	LinkName      string
	LogInConsole  bool
}
