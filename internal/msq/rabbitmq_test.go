package msq

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/captainmango/coco-cron-parser/internal/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

var rabbitMQHost string = "localhost:5672/"

func TestMain(m *testing.M) {
	os.Exit(runTests(m))
}

func runTests(m *testing.M) int {
	f, err := os.ReadFile(utils.BasePath("docker-compose.test.yml"))
	if err != nil {
		log.Printf("failed to create compose: %s", err)
		return 1
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stack, err := compose.NewDockerComposeWith(
		compose.WithStackReaders(
			strings.NewReader(string(f)),
		),
	)

	if err != nil {
		log.Printf("failed to create compose testing stack: %s", err)
		return 1
	}

	err = stack.
		WaitForService(
			"rabbitmq",
			wait.NewHTTPStrategy("/").
				WithPort("15672/tcp").
				WithStartupTimeout(30*time.Second),
		).
		Up(ctx, compose.Wait(true))

	if err != nil {
		log.Printf("failed to start compose testing stack: %s", err)
		return 1
	}

	defer func() {
		err = stack.Down(
			context.Background(),
			compose.RemoveOrphans(true),
			compose.RemoveVolumes(true),
		)

		if err != nil {
			log.Printf("failed to stop compose testing stack: %s", err)
		}
	}()

	return m.Run()
}

func Test_ItStartsCorrectly(t *testing.T) {
	rbmq, err := NewRabbitMQHandler(
		WithConnStr(rabbitMQHost),
	)

	require.NoError(t, err)
	require.NotNil(t, rbmq)
}

func Test_ItPublishesAndConsumesMessages(t *testing.T) {
	ready := make(chan struct{})

	rbmq, err := NewRabbitMQHandler(
		WithConnStr(rabbitMQHost),
		WithOnReadyFn(func() { close(ready) }),
	)

	require.NoError(t, err)
	require.NotNil(t, rbmq)
	defer rbmq.Close()

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	var receivedRoomId string
	received := make(chan struct{})
	defer close(received)

	go func() {
		consumeError := rbmq.ConsumeMessages(ctx, "coco_tasks.start_game", func(msg amqp.Delivery) error {
			var t struct {
				RoomId string `json:"room_id"`
			}

			json.NewDecoder(bytes.NewReader(msg.Body)).Decode(&t)
			receivedRoomId = t.RoomId

			msg.Ack(false)
			received <- struct{}{}

			return nil
		})

		require.NoError(t, consumeError)
	}()

	select {
	case <-ready:
	case <-time.After(5 * time.Second):
		log.Fatal("timed out waiting for topic to be available")
	}

	payload := StartGamePayload{
		RoomId: "23",
	}

	msgPayload, err := json.Marshal(payload)
	require.NoError(t, err)

	publishError := rbmq.PushMessage(payload.GetRoutingKey(), string(msgPayload))
	require.NoError(t, publishError)

	select {
	case <-received:
		// message was consumed. Channel closed
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for message to be consumed")
	}

	assert.Equal(t, "23", receivedRoomId)
}
