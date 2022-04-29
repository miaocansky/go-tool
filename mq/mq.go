package mq

import "github.com/streadway/amqp"

type MqEr interface {
	SendMessage(message Message, reliable bool) bool
	SendDelayMessage(message Message, reliable bool) bool
	Consumer(Receiver)
	Close()
}

// 定义接收者接口
type Receiver interface {
	Consumer(CustomerMsg) error
	FailAction(error, CustomerMsg) error
}

// 消息体：DelayTime 仅在 SendDelayMessage 方法有效
type Message struct {
	DelayTime int // desc:延迟时间(秒)
	Body      string
	RetryNum  int64 //第几次尝试了
}

// 消费者回调方法
type ConsumeFuc func(CustomerMsg) bool

type CustomerMsg struct {
	Body      []byte
	MessageId string
}

type Config struct {
	Username     string //用户名
	Password     string // 密码
	Host         string // 地址
	Port         string //端口
	Path         string // 根目录
	ExchangeName string // 交换器名称
	RouteKey     string // 路由名称
	QueueName    string // 队列名称
	RetryNum     int64  // 重试次数
}

type MessageQueue struct {
	Conn         *amqp.Connection // amqp链接对象
	Ch           *amqp.Channel    // channel对象
	ExchangeName string           // 交换器名称
	RouteKey     string           // 路由名称
	QueueName    string           // 队列名称
}
