package msq

import amqp "github.com/rabbitmq/amqp091-go"

type MessageQueueHandler interface {
	PushMessage(routingKey, body string) error
	PullMessages(routingKey string) (<-chan amqp.Delivery, error)
}
