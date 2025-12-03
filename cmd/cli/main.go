package main

import (
	"context"
	"log"
	"os"

	"github.com/captainmango/coco-cron-parser/internal/commands"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Commands: commands.Registry.All(),
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
        log.Fatal(err)
    }
}
