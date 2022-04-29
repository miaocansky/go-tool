package rabbitmq

import (
	"fmt"
	"github.com/gookit/goutil/dump"
	myMq "github.com/miaocansky/go-tool/mq"
	"github.com/streadway/amqp"
	"log"
	"strconv"
	"strings"
	"time"
)

// 消息体：DelayTime 仅在 SendDelayMessage 方法有效

type MessageQueue struct {
	conn         *amqp.Connection // amqp链接对象
	ch           *amqp.Channel    // channel对象
	ExchangeName string           // 交换器名称
	RouteKey     string           // 路由名称
	QueueName    string           // 队列名称
	RetryNum     int64            // 重试次数
	Url          string           //地址
	config       myMq.Config      // 配置
}

// NewRabbitMQ 新建 rabbitmq 实例
func NewRabbitMQ(config myMq.Config) (*MessageQueue, error) {
	messageQueue := &MessageQueue{
		ExchangeName: config.ExchangeName,
		RouteKey:     config.RouteKey,
		QueueName:    config.QueueName,
		RetryNum:     config.RetryNum,
	}

	url := fmt.Sprintf(
		"amqp://%s:%s@%s:%s%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		"/"+strings.TrimPrefix(config.Path, "/"),
	)
	messageQueue.Url = url
	//建立amqp链接
	conn, err := amqp.Dial(url)

	if err != nil {
		failOnError(err, "Failed to connect to RabbitMQ")

		return nil, err

	}

	//fmt.Sprintf(
	//	"amqp:%s:%s@%s:%s%s",
	//	config.Viper.GetString("rabbitmq.username"),
	//	config.Viper.GetString("rabbitmq.password"),
	//	config.Viper.GetString("rabbitmq.host"),
	//	config.Viper.GetString("rabbitmq.port"),
	//	"/"+strings.TrimPrefix(config.Viper.GetString("rabbitmq.vhost"), "/"),
	//)

	messageQueue.conn = conn

	// 建立channel通道
	ch, err := conn.Channel()
	if err != nil {
		failOnError(err, "Failed to open a channel")
		return nil, err

	}

	messageQueue.ch = ch

	// 声明exchange交换器
	err = messageQueue.declareExchange(config.ExchangeName, nil)
	if err != nil {
		return nil, err
	}
	return messageQueue, nil

}

// SendMessage 发送普通消息
func (mq *MessageQueue) SendMessage(message myMq.Message, reliable bool) bool {
	confirms := make(chan amqp.Confirmation, 1)
	defer close(confirms)
	if reliable {
		mq.ch.Confirm(false)
		confirms = mq.ch.NotifyPublish(confirms)
	}

	err := mq.ch.Publish(
		mq.ExchangeName, // exchange
		mq.RouteKey,     // route key
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(message.Body),
		},
	)
	if err != nil {
		failOnError(err, "send common msg err")
		return false
	}

	isOk := true
	if reliable {
		select {
		case value := <-confirms:
			if value.Ack {
				isOk = true
			} else {
				isOk = false
			}
			//fmt.Println("int:", value)
		case <-time.After(time.Second * 10):
			isOk = false
			fmt.Errorf("expected to close confirms on Channel.NotifyPublish chan after Connection.Close")
		}

	} else {

	}
	return isOk

}

/**
 * 消息重试发送
 */
func sendRetryMessage(msgBody []byte, retryNum int64, parentConfig myMq.Config) bool {
	retryMq, err := NewRabbitMQ(parentConfig)
	defer retryMq.Close()
	if err != nil {
		return false
	}
	return retryMq.sendRetryMessage(msgBody, retryNum, true)
}

/**
 * 发送重试消息
 */
func (mq *MessageQueue) sendRetryMessage(msgBody []byte, retryNum int64, reliable bool) bool {
	delayTime := 20
	delayQueueName := mq.QueueName + "__retry:" + strconv.FormatInt(mq.RetryNum, 10) + "_delay:" + strconv.Itoa(delayTime)
	delayRouteKey := mq.RouteKey + "__retry:" + strconv.FormatInt(mq.RetryNum, 10) + "_delay:" + strconv.Itoa(delayTime)
	confirms := make(chan amqp.Confirmation, 1)
	//defer close(confirms)
	if reliable {
		mq.ch.Confirm(false)
		confirms = mq.ch.NotifyPublish(confirms)
	}
	// 定义延迟队列(死信队列)
	dq := mq.declareQueue(
		delayQueueName,
		amqp.Table{
			"x-dead-letter-exchange":    mq.ExchangeName, // 指定死信交换机
			"x-dead-letter-routing-key": mq.RouteKey,     // 指定死信routing-key
		},
	)

	// 延迟队列绑定到exchange
	mq.bindQueue(dq.Name, delayRouteKey, mq.ExchangeName)
	retryNum = retryNum + 1
	headers := amqp.Table{
		"retry_num": retryNum,
	}

	// 发送消息，将消息发送到延迟队列，到期后自动路由到正常队列中
	err := mq.ch.Publish(
		mq.ExchangeName,
		delayRouteKey,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         msgBody,
			Expiration:   strconv.Itoa(delayTime * 1000),
			Headers:      headers,
		},
	)
	if err != nil {
		failOnError(err, "send common msg err")
		return false
	}

	isOk := true
	if reliable {
		select {
		case value := <-confirms:
			if value.Ack {
				isOk = true
			} else {
				isOk = false
			}
			//fmt.Println("int:", value)
		case <-time.After(time.Second * 10):
			isOk = false
			fmt.Errorf("expected to close confirms on Channel.NotifyPublish chan after Connection.Close")
		}
	} else {
	}
	return isOk
}

// SendDelayMessage 发送延迟消息
func (mq *MessageQueue) SendDelayMessage(message myMq.Message, reliable bool) bool {
	delayQueueName := mq.QueueName + "_delay:" + strconv.Itoa(message.DelayTime)
	delayRouteKey := mq.RouteKey + "_delay:" + strconv.Itoa(message.DelayTime)

	confirms := make(chan amqp.Confirmation, 1)
	//defer close(confirms)
	if reliable {
		//开启等待确认机制
		mq.ch.Confirm(false)
		confirms = mq.ch.NotifyPublish(confirms)
	}
	// 定义延迟队列(死信队列)
	dq := mq.declareQueue(
		delayQueueName,
		amqp.Table{
			"x-dead-letter-exchange":    mq.ExchangeName, // 指定死信交换机
			"x-dead-letter-routing-key": mq.RouteKey,     // 指定死信routing-key
		},
	)

	// 延迟队列绑定到exchange
	mq.bindQueue(dq.Name, delayRouteKey, mq.ExchangeName)

	//
	headers := amqp.Table{
		"retry_num": 0,
	}

	// 发送消息，将消息发送到延迟队列，到期后自动路由到正常队列中
	err := mq.ch.Publish(
		mq.ExchangeName,
		delayRouteKey,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(message.Body),
			Expiration:   strconv.Itoa(message.DelayTime * 1000),
			Headers:      headers,
		},
	)
	if err != nil {
		failOnError(err, "send common msg err")
		return false
	}

	isOk := true
	if reliable {
		select {
		case value := <-confirms:
			if value.Ack {
				isOk = true
			} else {
				isOk = false
			}
			//fmt.Println("int:", value)
		case <-time.After(time.Second * 10):
			isOk = false
			fmt.Errorf("expected to close confirms on Channel.NotifyPublish chan after Connection.Close")
		}
	} else {
	}
	return isOk

}

// Consume 获取消费消息
func (mq *MessageQueue) Consumer(receiver myMq.Receiver) {
	// 声明队列
	q := mq.declareQueue(mq.QueueName, nil)

	// 队列绑定到exchange
	mq.bindQueue(q.Name, mq.RouteKey, mq.ExchangeName)

	// 设置Qos
	err := mq.ch.Qos(1, 0, false)
	failOnError(err, "Failed to set QoS")

	// 监听消息
	msgs, err := mq.ch.Consume(
		q.Name, // queue name,
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)
	//注册在主进程，不需要阻塞
	go func() {
		for d := range msgs {
			//fmt.Println(d)
			//dump.Println(d.Headers)
			retryNum, ok := d.Headers["retry_num"].(int64)
			if !ok {
				retryNum = 0
			}
			s := string(d.Body)
			dump.Println(s + "=>" + strconv.FormatInt(retryNum, 10))
			msg := myMq.CustomerMsg{MessageId: d.MessageId, Body: d.Body}
			err := receiver.Consumer(msg)
			if err == nil {
				d.Ack(false)
			} else {
				if retryNum < mq.RetryNum {
					//重试
					retryOk := sendRetryMessage(d.Body, retryNum, mq.config)
					if !retryOk {
						d.Nack(false, true)
					} else {
						d.Ack(false)
					}

				} else {
					// 在队列中移除 并且记录到数据库中  重试失败处理
					receiver.FailAction(err, msg)
					d.Nack(false, false)

				}

			}
		}
	}()
	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

// Close 关闭链接
func (mq *MessageQueue) Close() {
	if mq.ch != nil {
		mq.ch.Close()
	}
	if mq.conn != nil {
		mq.conn.Close()
	}

}

// declareQueue 定义队列
func (mq *MessageQueue) declareQueue(name string, args amqp.Table) amqp.Queue {
	q, err := mq.ch.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		args,
	)
	failOnError(err, "Failed to declare a delay_queue")

	return q
}

// declareQueue 定义交换器
func (mq *MessageQueue) declareExchange(exchange string, args amqp.Table) error {
	err := mq.ch.ExchangeDeclare(
		exchange,
		"direct",
		true,
		false,
		false,
		false,
		args,
	)
	if err != nil {
		failOnError(err, "Failed to declare an exchange")
	}
	return err
}

// bindQueue 绑定队列
func (mq *MessageQueue) bindQueue(queue, routekey, exchange string) {
	err := mq.ch.QueueBind(
		queue,
		routekey,
		exchange,
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")
}

// failOnError 错误处理
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s : %s", msg, err)
	}
}
