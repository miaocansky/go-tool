package dto

type LoggerData struct {
	Key  string
	Data interface{}
}

func LoggerMsg(key string, data interface{}) LoggerData {
	return LoggerData{
		Key:  key,
		Data: data,
	}
}
