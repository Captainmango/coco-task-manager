package main

import (
	"fmt"
	"os"

	"github.com/captainmango/coco-cron-parser/internal/parser"
)

func main() {
	cronArgs := os.Args[1:]

	p, err := parser.NewParser(cronArgs)
	if err != nil {
		fmt.Print(err)
	}
	cron, _ := p.Parse()
	fmt.Println(cron)
}
