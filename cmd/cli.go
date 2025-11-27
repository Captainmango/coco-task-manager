package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/captainmango/coco-cron-parser/internal/data"
	"github.com/captainmango/coco-cron-parser/internal/parser"
)

func main() {
	cronArgs := strings.Join(os.Args[1:], " ")

	p, err := parser.NewParser(
		parser.WithInput(cronArgs, true),
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	cron, _ := p.Parse()
	cron.PrintingMode = data.POSSIBLE_VALUES
	fmt.Printf("%s\n", cron)
}
