package util

import (
	"os"

	pb "core/jobrequest"

	"github.com/golang/protobuf/proto"
	_ "github.com/joho/godotenv/autoload"
	"github.com/streadway/amqp"
)

type AmqpJobData = pb.JobRequest

func AmqpPublish(queue string, data *pb.JobRequest) (err error) {
	conn, err := amqp.Dial(os.Getenv("MQ_URL"))
	if ok := FailOnError(err, "Failed to connect to RabbitMQ"); ok {
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if ok := FailOnError(err, "Failed to open a channel"); ok {
		return
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		queue,    //name
		"fanout", //type
		true,     //durable
		false,    //auto-deleted
		false,    //internal
		false,    //no-wait
		nil,      //arguments
	)
	if ok := FailOnError(err, "Failed to declare a exchange"); ok {
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
	if ok := FailOnError(err, "Failed to declare a queue"); ok {
		return
	}

	body, _ := proto.Marshal(data)
	err = ch.Publish(
		q.Name, // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	FailOnError(err, "Failed to publish")
	return
}

func AmqpPublish2(mq_url string, queue string, data *pb.JobRequest) (err error) {
	if mq_url == "" {
		if os.Getenv("APP_ENV") == "production" {
			mq_url = "amqp://root:xxxxx@mq.prod.com:5672/"
		} else {
			mq_url = "amqp://guest:yyyyy@127.0.0.1:5672/"
		}
	}
	conn, err := amqp.Dial(mq_url)
	if ok := FailOnError(err, "Failed to connect to RabbitMQ"); ok {
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if ok := FailOnError(err, "Failed to open a channel"); ok {
		return
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		queue,    //name
		"fanout", //type
		true,     //durable
		false,    //auto-deleted
		false,    //internal
		false,    //no-wait
		nil,      //arguments
	)
	if ok := FailOnError(err, "Failed to declare a exchange"); ok {
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
	if ok := FailOnError(err, "Failed to declare a queue"); ok {
		return
	}

	body, _ := proto.Marshal(data)
	err = ch.Publish(
		q.Name, // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	FailOnError(err, "Failed to publish")
	return
}
