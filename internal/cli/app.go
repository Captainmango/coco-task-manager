package coco_cli

import (
	"log/slog"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/captainmango/coco-cron-parser/internal/config"
)

func CreateCLI() *cli.Command {
	config.BootstrapConfig(
		config.WithDotEnv(),
	)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	return &cli.Command{
		Commands: CommandRegistry.All(),
	}
}
