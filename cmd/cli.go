package main

import (
	"fmt"
	"os"

	"github.com/captainmango/coco-cron-parser/internal/parser"
)

func main() {
	cronArgs := os.Args[1]

	p := parser.NewParser(cronArgs)
	cron, _ := p.Parse()
	fmt.Println(cron)
}
