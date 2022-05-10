package log

var logManager = NewLog()

func Register(driverName string, driver Logger) *Log {
	logManager.Register(driverName, driver)
	return logManager
}

func DefaultDriver() Logger {
	return logManager.DefaultDriver()
}
