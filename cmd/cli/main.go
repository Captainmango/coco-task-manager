package main

import (
	"context"
	"log"
	"os"

	coco_cli "github.com/captainmango/coco-cron-parser/internal/cli"
)

func main() {
	cmd := coco_cli.CreateCLI()

	if err := cmd.Run(context.Background(), os.Args); err != nil {
        log.Fatal(err)
    }
}
