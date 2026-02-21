package msq

import (
	"context"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/captainmango/coco-cron-parser/internal/utils"
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
