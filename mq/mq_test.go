package mq

import (
	"errors"
	"fmt"
	"github.com/miaocansky/go-tool/mq/rabbitmq"
	"strconv"
	"testing"
	"time"
)

type RecvPro struct {
}

//// 实现消费者 消费消息失败 自动进入延时尝试  尝试3次之后入库db
/*
返回值 error 为nil  则表示该消息消费成功
否则消息会进入ttl延时队列  重复尝试消费3次
3次后消息如果还是失败 消息就执行失败  进入告警 FailAction
*/
func (t *RecvPro) Consumer(msg CustomerMsg) error {
	time.Sleep(time.Second * 1)
	//return errors.New("顶顶顶顶")
	fmt.Println(string(msg.Body))
	//time.Sleep(1*time.Second)
	return errors.New("消费异常")
	//return nil
}

//消息已经消费3次 失败了 请进行处理
/*
如果消息 消费多次后还是失败 进行入库记录
*/
func (t *RecvPro) FailAction(err error, msg CustomerMsg) error {
	fmt.Println(string(msg.Body))
	fmt.Println(err)
	fmt.Println("多次消费失败 入库")
	return nil
}

func TestPrduct(t *testing.T) {
	for i := 0; i < 8; i++ {
		msg := "这是一条普通消息" + strconv.Itoa(i)
		SendDelayMessage(msg)
		time.Sleep(1 * time.Second)
		//fmt.Println(i)
	}

}

func ConsumeDelay() {

	config := Config{
		Username:     "jenkin",
		Password:     "123456",
		Host:         "127.0.0.1",
		Port:         "5672",
		Path:         "/",
		ExchangeName: "exchange_2",
		RouteKey:     "route_2",
		QueueName:    "queue_2",
		RetryNum:     0,
	}

	rabbit, err := rabbitmq.NewRabbitMQ(config)
	defer rabbit.Close()
	if err != nil {

	}
	processTask := &RecvPro{}
	// 执行消费
	rabbit.Consumer(processTask)

}

func SendMessage() {
	config := Config{
		Username:     "jenkin",
		Password:     "123456",
		Host:         "127.0.0.1",
		Port:         "5672",
		Path:         "/",
		ExchangeName: "exchange",
		RouteKey:     "route",
		QueueName:    "queue",
		RetryNum:     0,
	}
	rabbit, err := rabbitmq.NewRabbitMQ(config)
	defer rabbit.Close()
	if err != nil {

	}
	message := rabbit.SendMessage(Message{Body: "这是一条普通消息"}, true)
	fmt.Println(message)

}

func SendDelayMessage(msg string) {

	config := Config{
		Username:     "jenkin",
		Password:     "123456",
		Host:         "127.0.0.1",
		Port:         "5672",
		Path:         "/",
		ExchangeName: "exchange_2",
		RouteKey:     "route_2",
		QueueName:    "queue_2",
		RetryNum:     0,
	}

	rabbit, err := rabbitmq.NewRabbitMQ(config)
	defer rabbit.Close()
	if err != nil {

	}
	rabbit.SendDelayMessage(Message{Body: msg, DelayTime: 5}, true)
	//fmt.Println(message)

}
