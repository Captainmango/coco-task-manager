package main

import (
	"context"
	"log"
	"os"

	_ "github.com/captainmango/coco-cron-parser/internal/resources"
	"github.com/captainmango/coco-cron-parser/internal/cli"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Commands: coco_cli.Registry.All(),
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
        log.Fatal(err)
    }
}
