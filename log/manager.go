package log

import "github.com/miaocansky/go-tool/log/dto"

type Log struct {
	defaultName string
	drivers     Logger
}

func Ang(key string, value interface{}) dto.LoggerData {
	return dto.LoggerData{
		Key:   key,
		Value: value,
	}
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
