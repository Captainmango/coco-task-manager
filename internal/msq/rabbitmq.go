package msq

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const exchangeName = "coco_tasks"

type RabbitMQOptFn func(*RabbitMQHandler)

func WithConnStr(hostname string) RabbitMQOptFn {
	return func(rm *RabbitMQHandler) {
		rm.host = fmt.Sprintf("amqp://%s", hostname)
	}
}

func WithBufferSize(bufferSize int) RabbitMQOptFn {
	return func(rm *RabbitMQHandler) {
		rm.bufferSize = bufferSize
	}
}

func NewRabbitMQHandler(opts ...RabbitMQOptFn) *RabbitMQHandler {
	rbmq := &RabbitMQHandler{}

	for _, opt := range opts {
		opt(rbmq)
	}

	if rbmq.bufferSize == 0 {
		rbmq.bufferSize = 50 // Sane default of 50 messages in buffer
	}

	return rbmq
}

type RabbitMQHandler struct {
	host       string
	bufferSize int
}

func (rbmq *RabbitMQHandler) PushMessage(routingKey, body string) error {
	return rbmq.handleConnection(func(ch *amqp.Channel) error {
		err := exchangeDefinition(ch)

		if err != nil {
			slog.Error(err.Error())
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = ch.PublishWithContext(
			ctx,
			exchangeName, // exchange
			routingKey,   // routing key
			false,        // mandatory
			false,        // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			},
		)

		if err != nil {
			slog.Error(err.Error())
			return err
		}

		slog.Info("Pushed message", slog.String("body", body))

		return nil
	})
}

func (rbmq *RabbitMQHandler) PullMessages(routingKey string) (<-chan amqp.Delivery, error) {
	msgs := make(<-chan amqp.Delivery, rbmq.bufferSize)
	err := rbmq.handleConnection(func(ch *amqp.Channel) error {
		err := exchangeDefinition(ch)
		if err != nil {
			slog.Error(err.Error())
			return err
		}

		q, err := ch.QueueDeclare(
			"",    // name
			false, // durable
			false, // delete when unused
			true,  // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		err = ch.Qos(rbmq.bufferSize, 0, false)
		if err != nil {
			slog.Error(err.Error())
			return err
		}

		err = ch.QueueBind(
			q.Name,       // queue name
			routingKey,   // routing key
			exchangeName, // exchange
			false,        // no wait
			nil,          // arguments
		)
		if err != nil {
			slog.Error(err.Error())
			return err
		}

		msgs, err = ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto ack
			false,  // exclusive
			false,  // no local
			false,  // no wait
			nil,    // args
		)
		if err != nil {
			slog.Error(err.Error())
			return err
		}

		return nil
	})

	return msgs, err
}

func (rbmq *RabbitMQHandler) handleConnection(fn func(ch *amqp.Channel) error) error {
	conn, err := amqp.Dial(rbmq.host)

	if err != nil {
		slog.Error(err.Error())
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	defer ch.Close()

	if err = fn(ch); err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}

func exchangeDefinition(ch *amqp.Channel) error {
	err := ch.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

	return err
}
