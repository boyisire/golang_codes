package base

import (
	"fmt"
	"os"
	"strconv"
	"time"

	pb "core/jobrequest"
	"core/log"

	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
)

// RabbitMq 消息类型
type RabbitMq struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      amqp.Queue
	Exchange   string
}

// Close close connection
func (me *RabbitMq) Close() {
	me.Connection.Close()
	me.Channel.Close()
}

func connAmqp(url string) *amqp.Connection {
	conn, err := amqp.Dial(url)
	if err != nil {
		time.Sleep(time.Duration(retries) * time.Second)
		msg := fmt.Sprintf("RabbitMQ Connection Fail:: err=[%s], retries=[%d], url=[%s]", err.Error(), retries, url)
		log.Error(msg, log.Data{
			"url":       url,
			"log_to_es": false,
		})
		retries += 2
		connAmqp(url)
	}
	return conn
}

// NewRabbitMq 初始化rabbitmq
func NewRabbitMq(url, queue string) (mq *RabbitMq) {
	mq = &RabbitMq{}
	conn := connAmqp(url)
	if conn == nil {
		msg := fmt.Sprintf("RabbitMQ NewRabbitMq Fail,url=[%s], queue=[%s]", url, queue)
		log.Error(msg, log.Data{
			"url":       url,
			"queue":     queue,
			"log_to_es": false,
		})
		return
	}
	mq.Connection = conn

	ch, err := conn.Channel()
	if err != nil {
		msg := fmt.Sprintf("RabbitMQ Open Channel Fail:: err=[%s]", err.Error())
		log.Error(msg, log.Data{
			"url":   url,
			"queue": queue,
		})
		return
	}
	mq.Channel = ch

	err = ch.ExchangeDeclare(
		queue,    // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		msg := fmt.Sprintf("rabbitMq Declare Exchange Fail:: err=[%s]", err.Error())
		log.Error(msg, log.Data{
			"url":   url,
			"queue": queue,
		})
		return
	}

	q, err := ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		msg := fmt.Sprintf("RabbitMQ Declare Queue Fail:: err=[%s]", err.Error())
		log.Error(msg, log.Data{
			"url":   url,
			"queue": queue,
		})
		return
	}
	mq.Queue = q

	err = ch.QueueBind(
		q.Name, // queue name
		q.Name, // routing key
		q.Name, // exchange
		false,
		nil,
	)
	if err != nil {
		msg := fmt.Sprintf("RabbitMQ Bind Queue Fail:: err=[%s]", err.Error())
		log.Error(msg, log.Data{
			"url":   url,
			"queue": queue,
		})
		return
	}
	mq.Exchange = q.Name

	return
}

// Consume 处理消息
func (me *RabbitMq) Consume(job IJob, async bool) {
	defer func() {
		me.Close()
	}()
	var autoAck = true
	//异步消费
	if async == true {
		autoAck = false
	}
	//限制协程数量
	mqChannelNum := os.Getenv("MQ_CHANNEL_NUM")
	if mqChannelNum == "" {
		mqChannelNum = "1000"
	}
	mqN, _ := strconv.Atoi(mqChannelNum)
	mqNch := make(chan int, mqN)

	msgs, err := me.Channel.Consume(
		me.Queue.Name, // queue
		"",            // consumer
		autoAck,       // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		msg := fmt.Sprintf("RabbitMQ Register Consumer Fail:: err=[%s]", err.Error())
		log.Error(msg, log.Data{
			"autoAck": autoAck,
			"queue":   me.Queue.Name,
		})
		return
	}

	log.New().WithFields(log.Fields{
		"exchange": me.Exchange,
		"queue":    me.Queue.Name,
	}).Info(" [*] Waiting for msgs. To exit press CTRL+C")

	//处理消息
	for d := range msgs {
		//protobuf解码
		preq := &pb.JobRequest{}

		err := proto.Unmarshal(d.Body, preq)
		if err != nil {
			log.Error(fmt.Sprintf("Consume:Proto unmarshaling error. err=%s", err.Error()))
		}
		if async == true {
			mqNch <- 1
			go func(d amqp.Delivery, ch chan int) {
				job.Handle(preq)
				d.Ack(false)
				<-ch
			}(d, mqNch)
		} else {
			job.Handle(preq)
		}
	}
}
