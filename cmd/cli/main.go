package main

import (
	"context"
	"log"
	"os"

	coco_cli "github.com/captainmango/coco-cron-parser/internal/cli"
	"github.com/captainmango/coco-cron-parser/internal/config"
)

func main() {
	config.BootstrapConfig(
		config.WithDotEnv(),
	)

	cmd := coco_cli.CreateCLI()

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
