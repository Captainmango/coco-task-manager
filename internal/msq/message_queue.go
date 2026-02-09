package msq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumeMessageFn func(msg amqp.Delivery) error

type AdvancedMessageQueueHandler interface {
	PushMessage(routingKey, body string) error
	ConsumeMessages(ctx context.Context, routingKey string, fn ConsumeMessageFn) error
}
