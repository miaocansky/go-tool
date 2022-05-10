package log

type Log struct {
	defaultName string
	drivers     Logger
}

func NewLog() *Log {
	return &Log{}
}
func (log *Log) SetDefaultName(defaultName string) {
	log.defaultName = defaultName
}
func (log *Log) Register(defaultName string, logger Logger) {
	log.defaultName = defaultName
	log.drivers = logger
}
func (log *Log) DefaultDriver() Logger {
	return log.drivers
}
