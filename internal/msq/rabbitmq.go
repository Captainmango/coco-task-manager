package msq

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const exchangeName = "coco_tasks"

type RabbitMQOptFn func(*RabbitMQHandler) error

func WithConnStr(hostname string) RabbitMQOptFn {
	return func(rm *RabbitMQHandler) error {
		rm.host = fmt.Sprintf("amqp://%s", hostname)
		return nil
	}
}

func WithBufferSize(bufferSize int) RabbitMQOptFn {
	return func(rm *RabbitMQHandler) error {
		rm.bufferSize = bufferSize
		return nil
	}
}

func NewRabbitMQHandler(opts ...RabbitMQOptFn) (*RabbitMQHandler, error) {
	rbmq := &RabbitMQHandler{}

	for _, opt := range opts {
		if err := opt(rbmq); err != nil {
			return nil, err
		}
	}

	if rbmq.bufferSize == 0 {
		rbmq.bufferSize = 50 // Sane default of 50 messages in buffer
	}

	conn, err := amqp.Dial(rbmq.host)

	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	rbmq.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	rbmq.ch = ch

	return rbmq, nil
}

func (rbmq *RabbitMQHandler) Close() error {
	rbmq.mu.Lock()
	chann := rbmq.ch
	rbmq.mu.Unlock()

	err := chann.Close()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	err = rbmq.conn.Close()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}

type RabbitMQHandler struct {
	host       string
	bufferSize int
	conn       *amqp.Connection
	ch         *amqp.Channel
	mu         sync.Mutex
}

func (rbmq *RabbitMQHandler) PushMessage(routingKey, body string) error {
	rbmq.mu.Lock()
	chann := rbmq.ch
	rbmq.mu.Unlock()
	if chann == nil {
		return errors.New("rabbitmq channel is closed")
	}

	if err := exchangeDefinition(chann); err != nil {
		slog.Error(err.Error())
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := chann.PublishWithContext(
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
}

func (rbmq *RabbitMQHandler) ConsumeMessages(ctx context.Context, routingKey string, fn ConsumeMessageFn) error {
	rbmq.mu.Lock()
	chann := rbmq.ch
	rbmq.mu.Unlock()
	if chann == nil {
		return errors.New("rabbitmq channel is closed")
	}

	err := exchangeDefinition(chann)

	if err != nil {
		slog.Error(err.Error())
		return err
	}

	q, err := chann.QueueDeclare(
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

	if err = chann.Qos(rbmq.bufferSize, 0, false); err != nil {
		slog.Error(err.Error())
		return err
	}

	err = chann.QueueBind(
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

	msgs, err := chann.Consume(
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

	slog.Info("Consuming messages",
		slog.String("topic", routingKey),
	)

	for {
		select {
		case m := <-msgs:
			if err = fn(m); err != nil {
				return err
			}
		case <-ctx.Done():
			slog.Info("Stopping message consumption")
			return nil
		}
	}
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
