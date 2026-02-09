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

type RabbitMQHandler struct {
	host       string
	bufferSize int
	conn       *amqp.Connection
	chPool     chan *amqp.Channel
	mu         sync.RWMutex
	closed     bool
	poolSize   int
	timeout    time.Duration
}

type RabbitMQOptFn func(*RabbitMQHandler) error

func NewRabbitMQHandler(opts ...RabbitMQOptFn) (*RabbitMQHandler, error) {
	rbmq := &RabbitMQHandler{
		poolSize:   10,
		bufferSize: 50,
		timeout:    5 * time.Second,
	}

	for _, opt := range opts {
		if err := opt(rbmq); err != nil {
			return nil, err
		}
	}

	conn, err := amqp.Dial(rbmq.host)

	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	rbmq.conn = conn

	rbmq.chPool = make(chan *amqp.Channel, rbmq.bufferSize)
	for i := 0; i < rbmq.bufferSize; i++ {
		ch, err := conn.Channel()
		if err != nil {
			rbmq.Close()
			slog.Error(err.Error())
			return nil, fmt.Errorf("unable to create Channel: %w", err)
		}

		rbmq.chPool <- ch
	}

	return rbmq, nil
}

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

func WithPoolSize(poolSize int) RabbitMQOptFn {
	return func(rm *RabbitMQHandler) error {
		rm.poolSize = poolSize
		return nil
	}
}

func WithConnectTimeout(timeout time.Duration) RabbitMQOptFn {
	return func(rm *RabbitMQHandler) error {
		rm.timeout = timeout
		return nil
	}
}

func (rbmq *RabbitMQHandler) getChannel() (*amqp.Channel, error) {
	rbmq.mu.RLock()
	if rbmq.closed {
		rbmq.mu.RUnlock()
		return nil, errors.New("handler closed")
	}
	rbmq.mu.RUnlock()

	select {
	case ch := <-rbmq.chPool:
		return ch, nil
	case <-time.After(rbmq.timeout):
		return nil, errors.New("timeout waiting for channel")
	}
}

func (rbmq *RabbitMQHandler) returnChannel(ch *amqp.Channel) {
	rbmq.mu.RLock()
	defer rbmq.mu.RUnlock()

	if rbmq.closed {
		ch.Close()
		return
	}

	select {
	case rbmq.chPool <- ch:
	default:
		// Unreachable, but just in case...
		ch.Close()
	}
}

func (rbmq *RabbitMQHandler) Close() error {
	rbmq.mu.Lock()
	if rbmq.closed {
		rbmq.mu.Unlock()
		return nil
	}
	rbmq.closed = true
	rbmq.mu.Unlock()

	close(rbmq.chPool)
	for ch := range rbmq.chPool {
		if err := ch.Close(); err != nil {
			slog.Error(err.Error())
		}
	}

	if err := rbmq.conn.Close(); err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}

func (rbmq *RabbitMQHandler) PushMessage(routingKey, body string) error {
	ch, err := rbmq.getChannel()

	if err != nil {
		return err
	}
	defer rbmq.returnChannel(ch)

	if err := exchangeDefinition(ch); err != nil {
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
}

func (rbmq *RabbitMQHandler) ConsumeMessages(ctx context.Context, routingKey string, fn ConsumeMessageFn) error {
	rbmq.mu.RLock()
	if rbmq.closed {
		rbmq.mu.Unlock()
		return errors.New("handler closed")
	}
	conn := rbmq.conn
	rbmq.mu.RUnlock()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	if err = exchangeDefinition(ch); err != nil {
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

	if err = ch.Qos(rbmq.bufferSize, 0, false); err != nil {
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

	msgs, err := ch.Consume(
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
